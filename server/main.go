package main

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"urjith.dev/algobattle/internal/bot"
	"urjith.dev/algobattle/internal/handlers"
	"urjith.dev/algobattle/pkg/services"
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
	defer db.Close()
	if err != nil {
		fmt.Printf("error creating firestore client: %v", err)
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.RecoveryWithWriter(os.Stdout))

	tiingo := services.NewTiingo(os.Getenv("TIINGO_TOKEN"))

	botworker := bot.NewBotWorker(db, tiingo)

	handlers.SetupRoutes(r, botworker)

	r.Run(":8080")
}
