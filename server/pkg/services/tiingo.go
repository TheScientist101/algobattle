// Package services provides external API integrations and data services
// for the AlgoBattle trading platform.
package services

import (
	"cmp"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
	"urjith.dev/algobattle/pkg/indicators"
	"urjith.dev/algobattle/pkg/models"
	"urjith.dev/algobattle/pkg/utils"
)

// Constants for Tiingo API configuration and caching
const (
	baseURL        = "https://api.tiingo.com" // Base URL for Tiingo API
	dataStart      = "1900-01-01"             // Start date for historical data
	dailyFreq      = "daily"                  // Frequency for historical data
	cacheFolder    = "./data"                 // Folder for caching data
	dailyCacheJSON = "dailycache.json"        // JSON cache filename
	dailyCacheGOB  = "dailycache.gob"         // GOB cache filename
)

// Tiingo is a client for the Tiingo API that provides stock market data.
// It manages a list of watched tickers, caches historical data, and
// calculates technical indicators.
type Tiingo struct {
	Token      string                 // API token for authentication
	tickers    *utils.TreeSet[string] // Set of watched ticker symbols
	DailyCache *models.History        // Cache of historical daily data
	Indicators []indicators.Indicator // Technical indicators to calculate
}

// NewTiingo creates a new Tiingo client with the provided API token.
// It initializes the ticker set, daily cache, and indicators list.
func NewTiingo(token string) *Tiingo {
	return &Tiingo{
		token,
		utils.NewTreeSet[string](cmp.Compare), // Create sorted set for tickers
		models.NewHistory(),                   // Initialize empty history
		make([]indicators.Indicator, 0),       // Initialize empty indicators list
	}
}

// AddTickers adds one or more ticker symbols to the watchlist.
// All tickers are converted to uppercase before being added.
func (t *Tiingo) AddTickers(newTickers ...string) {
	// Convert all tickers to uppercase
	for i, ticker := range newTickers {
		newTickers[i] = strings.ToUpper(ticker)
	}

	// Add tickers to the set
	t.tickers.Insert(newTickers...)
}

// LastPriceResponse represents the response from the Tiingo API for last price.
// This struct maps to the JSON response from the IEX endpoint.
type LastPriceResponse struct {
	Ticker   string  `json:"ticker"`   // Ticker symbol
	TngoLast float64 `json:"tngoLast"` // Latest price
}

// FetchCurrPrices fetches the current prices for all tickers in the watchlist.
// It makes a single API call to get prices for all tickers and returns a map
// of ticker symbols to their current prices.
func (t *Tiingo) FetchCurrPrices() map[string]float64 {
	tickers := t.tickers.AsSlice()
	tickersStr := strings.Join(tickers, ",")

	request, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/iex/?tickers=%s&token=%s",
			baseURL,
			tickersStr,
			t.Token,
		),
		nil)
	if err != nil {
		log.Println(err)
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			log.Println(tickers, "not found")
		}

		log.Println(response.Status+" when fetching ", tickers)
	}

	result := make([]LastPriceResponse, len(tickers))
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	prices := make(map[string]float64, len(tickers))
	for _, pair := range result {
		prices[pair.Ticker] = pair.TngoLast
	}

	return prices
}

// HistoricalDaily fetches historical daily data for a specific ticker.
// It retrieves data from the earliest available date and adds it to the daily cache.
// Returns an error if the API request fails or if the ticker is not found.
func (t *Tiingo) HistoricalDaily(ticker string) error {
	request, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(
			"%s/tiingo/daily/%s/prices?startDate=%s&resampleFreq=%s&format=%s&token=%s",
			baseURL,
			ticker,
			dataStart,
			dailyFreq,
			"json",
			t.Token,
		),
		nil,
	)

	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			log.Println(ticker, "not found")
			t.tickers.Remove(ticker)
		}

		return fmt.Errorf(response.Status + " when fetching " + ticker)
	}

	results := make([]models.PackedPeriod, 0, 365*5) // Pre-allocate 5 years of daily data
	if err = json.NewDecoder(response.Body).Decode(&results); err != nil {
		return err
	}

	t.DailyCache.AddData(results, ticker)

	return nil
}

// LoadData loads data from cache and downloads missing data for all tickers.
// It first tries to load from cache files, then downloads any missing ticker data.
// The useJSON parameter determines whether to use JSON or GOB format for loading.
func (t *Tiingo) LoadData(useJSON bool) error {
	if len(t.DailyCache.Rows) != 0 {
		log.Println("Warning := Overwriting DailyCache with file data")
	}

	err := t.LoadCaches(useJSON)
	if err != nil {
		return err
	}

	errs, _ := errgroup.WithContext(context.Background())

	log.Println("Downloading uncached tickers...")
	for ticker := range t.tickers.All() {
		if _, ok := t.DailyCache.Tickers[ticker]; !ok {
			// Use a closure to capture the ticker value correctly
			ticker := ticker // Create a new variable for the closure
			errs.Go(func() error {
				return t.HistoricalDaily(ticker)
			})
		}
	}

	err = errs.Wait()

	if err := t.SaveCaches(); err != nil {
		return err
	}

	return err
}

// DownloadAllTickers downloads data for all tickers
func (t *Tiingo) DownloadAllTickers() error {
	errs, _ := errgroup.WithContext(context.Background())

	for ticker := range t.tickers.All() {
		errs.Go(func() error {
			return t.HistoricalDaily(ticker)
		})
	}

	err := errs.Wait()

	if err := t.SaveCaches(); err != nil {
		return err
	}

	return err
}

// DownloadMissingTickers downloads data for tickers not in the cache
func (t *Tiingo) DownloadMissingTickers() error {
	errs, _ := errgroup.WithContext(context.Background())

	for ticker := range t.tickers.All() {
		if _, ok := t.DailyCache.Tickers[ticker]; !ok {
			errs.Go(func() error {
				return t.HistoricalDaily(ticker)
			})
		}
	}

	err := errs.Wait()

	if err := t.SaveCaches(); err != nil {
		return err
	}

	return err
}

// LoadCaches loads historical stock data caches from disk.
// If useJSON is true, it loads from the JSON cache file, otherwise from the GOB file.
// It creates the cache directory if it doesn't exist.
func (t *Tiingo) LoadCaches(useJSON bool) error {
	if useJSON {
		err := os.Mkdir(cacheFolder, 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		if _, err = os.Stat(filepath.Join(cacheFolder, dailyCacheJSON)); !errors.Is(err, os.ErrNotExist) {
			read, err := os.Open(filepath.Join(cacheFolder, dailyCacheJSON))
			if err == nil {
				err = json.NewDecoder(read).Decode(&t.DailyCache)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	}

	file, err := os.OpenFile(filepath.Join(cacheFolder, dailyCacheGOB), os.O_RDONLY, 0777)
	if err != nil {
		return err
	}

	packed := &models.PackedHistory{}
	err = gob.NewDecoder(file).Decode(packed)
	if err != nil {
		return err
	}

	t.DailyCache = packed.Unpack()

	return nil
}

// SaveCaches saves the daily cache to disk in both GOB and JSON formats.
// GOB format is used for efficient loading, while JSON is more portable.
// It creates the cache directory if it doesn't exist.
func (t *Tiingo) SaveCaches() error {
	err := os.Mkdir(cacheFolder, 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}

	packed := t.DailyCache.Pack()

	file, err := os.OpenFile(filepath.Join(cacheFolder, dailyCacheGOB), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println(err)
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(packed)
	if err != nil {
		log.Println(err)
	}

	marshalled, err := json.Marshal(packed)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(cacheFolder, dailyCacheJSON), marshalled, 0644)
	if err != nil {
		return err
	}

	return nil
}

// AddIndicator adds an indicator to the list
func (t *Tiingo) AddIndicator(indicator indicators.Indicator) {
	t.Indicators = append(t.Indicators, indicator)
}

// CalculateIndicators calculates all indicators for the daily cache
func (t *Tiingo) CalculateIndicators() error {
	log.Println("Calculating indicators...")

	indicators.CalculateIndicators(t.DailyCache, t.Indicators)

	return t.SaveCaches()
}
