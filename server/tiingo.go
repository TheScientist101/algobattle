package main

import (
	"cmp"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const baseURL = "https://api.tiingo.com"
const dataStart = "1900-01-01" // for Rows caching and downloading
const dailyFreq = "daily"      // for historical daily calls
const cacheFolder = "./data"
const dailyCacheJSON = "dailycache.json"
const dailyCacheGOB = "dailycache.gob"

type Tiingo struct {
	Token      string
	tickers    *TreeSet[string]
	DailyCache *History
	Indicators []Indicator
}

func NewTiingo(token string) *Tiingo {
	return &Tiingo{
		token,
		NewTreeSet[string](cmp.Compare),
		NewHistory(),
		make([]Indicator, 0),
	}
}

func (t *Tiingo) AddTickers(newTickers ...string) {
	for i, ticker := range newTickers {
		newTickers[i] = strings.ToUpper(ticker)
	}

	t.tickers.Insert(newTickers...)
}

type LastPriceResponse struct {
	Ticker   string  `json:"ticker"`
	TngoLast float64 `json:"tngoLast"`
}

func (t *Tiingo) fetchCurrPrices() map[string]float64 {
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
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			log.Println(tickers, "not found")
		}
		log.Fatal(response.Status+" when fetching ", tickers)
	}

	result := make([]LastPriceResponse, len(tickers))
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	mappings := make(map[string]float64, len(tickers))

	for _, pair := range result {
		mappings[pair.Ticker] = pair.TngoLast
	}

	return mappings
}

func (t *Tiingo) historicalDaily(ticker string) error {
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
		}

		t.tickers.Remove(ticker)

		return fmt.Errorf(response.Status + " when fetching " + ticker)
	}

	results := make([]PackedPeriod, 0, 365*5)
	if err = json.NewDecoder(response.Body).Decode(&results); err != nil {
		return err
	}

	t.DailyCache.addData(results, ticker)

	return nil
}

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
			errs.Go(func() error {
				return t.historicalDaily(ticker)
			})
		}
	}

	err = errs.Wait()

	if err := t.SaveCaches(); err != nil {
		return err
	}

	return err
}

func (t *Tiingo) DownloadAllTickers() error {
	errs, _ := errgroup.WithContext(context.Background())

	for ticker := range t.tickers.All() {
		errs.Go(func() error {
			return t.historicalDaily(ticker)
		})
	}

	err := errs.Wait()

	if err := t.SaveCaches(); err != nil {
		return err
	}

	return err
}

func (t *Tiingo) DownloadMissingTickers() error {
	errs, _ := errgroup.WithContext(context.Background())

	for ticker := range t.tickers.All() {
		if _, ok := t.DailyCache.Tickers[ticker]; !ok {
			errs.Go(func() error {
				return t.historicalDaily(ticker)
			})
		}
	}

	err := errs.Wait()

	if err := t.SaveCaches(); err != nil {
		return err
	}

	return err
}

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

	packed := &PackedHistory{}
	err = gob.NewDecoder(file).Decode(packed)
	if err != nil {
		return err
	}

	t.DailyCache = packed.Unpack()

	return nil
}

func (t *Tiingo) SaveCaches() error {
	err := os.Mkdir(cacheFolder, 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}

	file, err := os.OpenFile(filepath.Join(cacheFolder, dailyCacheGOB), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println(err)
	}

	packed := t.DailyCache.Pack()

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

func (t *Tiingo) AddIndicator(indicator Indicator) {
	t.Indicators = append(t.Indicators, indicator)
}

func (t *Tiingo) CalculateIndicators() error {
	log.Println("Calculating indicators...")

	t.DailyCache.CalculateIndicators(t.Indicators)

	return t.SaveCaches()
}
