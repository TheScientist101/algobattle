package main

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/olahol/melody"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file\n")
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS_FILE_PATH"))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	db, err := app.Firestore(ctx)
	if err != nil {
		fmt.Printf("error creating firestore client: %v", err)
	}

	r := gin.Default()
	m := melody.New()

	r.Use(gin.Logger())
	r.Use(gin.RecoveryWithWriter(os.Stdout))

	r.GET("/trading", func(c *gin.Context) {
		err := m.HandleRequest(c.Writer, c.Request)
		if err != nil {
			log.Printf("error establishing websocket connection: %v\n", err)
		}
	})

	httpRoutes := r.Group("/")

	botworker := NewBotWorker(db, NewTiingo(os.Getenv("TIINGO_TOKEN")))

	httpRoutes.Use(botworker.AuthHandler)

	httpRoutes.GET("/portfolio", botworker.GetPortfolio)
	httpRoutes.GET("/add_ticker", botworker.AddTicker)
	httpRoutes.POST("/transact", botworker.MakeTransaction, botworker.SavePortfolio)
	httpRoutes.GET("/stock_data", botworker.GetStockData)

	// TODO: Websockets
	//m.HandleMessage(botworker.TradingStream)
	//m.HandleSentMessage(func(s *melody.Session, bytes []byte) {
	//	ref, ok := s.Get("db_ref")
	//	if !ok {
	//		log.Println("db ref not found")
	//	} else if portfolio, ok := s.Get("bot"); ok {
	//		_, err = ref.(*firestore.DocumentRef).Set(context.Background(), portfolio.(Portfolio))
	//		if err != nil {
	//			log.Printf("error setting portfolio: %v\n", err)
	//		}
	//	} else {
	//		log.Println("error saving portfolio to firestore")
	//	}
	//})

	r.Run(":8080")
	defer db.Close()
}
