package handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"urjith.dev/algobattle/internal/bot"
)

// SetupRoutes sets up the HTTP routes
func SetupRoutes(r *gin.Engine, botWorker *bot.BotWorker) {
	httpRoutes := r.Group("/")
	httpRoutes.Use(botWorker.AuthHandler)

	httpRoutes.GET("/portfolio", botWorker.GetPortfolio)
	httpRoutes.GET("/add_ticker", botWorker.AddTicker)
	httpRoutes.POST("/transact", botWorker.MakeTransaction, botWorker.SavePortfolio)
	httpRoutes.GET("/daily_stock_data", botWorker.GetDailyStockData)
	httpRoutes.GET("/live_stock_data", botWorker.GetLiveStockData)
}

// DataPacket represents a data packet sent over WebSocket
type DataPacket struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// JSON converts the DataPacket to JSON
func (dp *DataPacket) JSON() []byte {
	b, err := json.Marshal(dp)
	if err != nil {
		panic(err)
	}

	return b
}

// ResultData represents a result message
type ResultData struct {
	Message string `json:"payload"`
	Success bool   `json:"success"`
}

// NewResultPacket creates a new result packet
func NewResultPacket(message string, success bool) *DataPacket {
	return &DataPacket{
		Type:    "result",
		Payload: &ResultData{message, success},
	}
}
