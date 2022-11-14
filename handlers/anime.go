package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AnimeTrackerr/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

// setup
func Setup(collectionRef *mongo.Collection) {
	collection = collectionRef
}

// handlers
func Default(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome To AnimeTracker :)")
}

// get list of anime
func GetAnime(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	opt := options.Find()
    limit, present := query["limit"]
	
	// check if query params "limit" is empty
    if !present || len(limit) == 0 {
        w.WriteHeader(400)
		fmt.Fprintf(w,"Error: Bad Request")
		return
    }

	lmt, _ := strconv.ParseInt(limit[0],10,64)

	// check if empty documents are being accessed
	if lmt <= 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w,"Error: Bad Request")
		return
	}

	offset, present := query["offset"]

	if present {
		oset, _ := strconv.ParseInt(offset[0],10,64)
	
		// check if offset is within range of docs
		length, _ := collection.CountDocuments(context.TODO(), bson.M{})
	
		if oset > length {
			w.WriteHeader(404)
			fmt.Fprintf(w,"Error: Documents trying to access do not exist")
			return
		}

		if oset < 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w,"Error: Bad Request")
			return
		}

		opt.SetSkip(oset)
	}
	
	opt.SetLimit(lmt)

	cursor, err := collection.Find(context.TODO(),bson.M{},opt)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			w.WriteHeader(404)
			fmt.Fprintf(w,"no documents found")
			return
		}
		panic(err)
	}

	var result []models.Anime
	if err = cursor.All(context.Background(), &result); err != nil {
  		log.Fatal(err)
	}

	// return json object through HTTP response
	json.NewEncoder(w).Encode(result) 
}

// search anime based on title
func SearchAnime(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	title, present := query["title"]

	if !present {
		w.WriteHeader(400)
		fmt.Fprintf(w,"Error: Bad Request")
		return
	}

	filter := bson.M {
		"title": bson.M {
			"$regex": title[0],
			"$options": "i",
		},
	}

	cursor, err := collection.Find(context.TODO(),filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			w.WriteHeader(404)
			fmt.Fprintf(w,"no documents found")
			return
		}
		panic(err)
	}

	var result []models.Anime
	if err = cursor.All(context.Background(), &result); err != nil {
  		log.Fatal(err)
	}

	// return json object through HTTP response
	json.NewEncoder(w).Encode(result) 
}