package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AnimeTrackerr/backend/handlers"
	"github.com/AnimeTrackerr/backend/server"
	"github.com/AnimeTrackerr/backend/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func routes(collection *mongo.Collection) {
	handlers.Setup(collection)

	http.HandleFunc("/", handlers.Default)
	http.HandleFunc("/getanime", handlers.GetAnime)
	http.HandleFunc("/searchanime",handlers.SearchAnime)
	http.ListenAndServe(os.Getenv("PORT"),nil)
}

func main() {
	// load env variables from file(local) or from cloud
	_ , err := os.OpenFile("config/.env",os.O_RDONLY,002)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("setting from cloud env")
	} else if err == nil{
		err := godotenv.Load("config/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		log.Fatal("Error opening File")
	}

	// connect to DB
	var cinfo utils.ClientInfo = server.ConnectDB(os.Getenv("URI"))
	var collection *mongo.Collection = server.GetCollection(cinfo.Client,os.Getenv("DB"),os.Getenv("COLLECTION"))

	defer cinfo.Client.Disconnect(cinfo.Ctx)

	// handle routes
	routes(collection)
}