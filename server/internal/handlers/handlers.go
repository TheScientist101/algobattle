// Package handlers provides HTTP route setup and WebSocket message handling
// for the AlgoBattle trading platform API.
package handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"urjith.dev/algobattle/internal/bot"
)

// SetupRoutes configures all HTTP routes for the application API.
// It groups routes under authentication middleware and maps each endpoint
// to its corresponding handler function in the BotWorker.
func SetupRoutes(r *gin.Engine, botWorker *bot.BotWorker) {
	httpRoutes := r.Group("/")
	httpRoutes.Use(botWorker.AuthHandler)

	httpRoutes.GET("/portfolio", botWorker.GetPortfolio)
	httpRoutes.GET("/add_ticker", botWorker.AddTicker)
	httpRoutes.POST("/transact", botWorker.MakeTransaction, botWorker.SavePortfolio)
	httpRoutes.GET("/daily_stock_data", botWorker.GetDailyStockData)
	httpRoutes.GET("/live_stock_data", botWorker.GetLiveStockData)
}

// DataPacket represents a data packet sent over WebSocket.
// It contains a type identifier and a payload that can be any type of data.
type DataPacket struct {
	Type    string `json:"type"`    // Type identifies the kind of data being sent
	Payload any    `json:"payload"` // Payload contains the actual data
}

// JSON converts the DataPacket to JSON byte array for transmission.
// It will panic if marshaling fails, which should only happen if the payload
// contains types that cannot be marshaled to JSON.
func (dp *DataPacket) JSON() []byte {
	b, err := json.Marshal(dp)
	if err != nil {
		panic(err)
	}

	return b
}

// ResultData represents a result message with success status.
// It is used to provide feedback for API operations.
type ResultData struct {
	Message string `json:"payload"` // Human-readable message describing the result
	Success bool   `json:"success"` // Indicates whether the operation was successful
}

// NewResultPacket creates a new result packet with the specified message and success status.
// This is a convenience function for creating standardized response packets.
func NewResultPacket(message string, success bool) *DataPacket {
	return &DataPacket{
		Type:    "result",
		Payload: &ResultData{message, success},
	}
}
