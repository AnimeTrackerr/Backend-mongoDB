package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AnimeTrackerr/backend/models"
	"github.com/AnimeTrackerr/backend/utils"
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
	if r.Method != "GET" {
		w.WriteHeader(405)
		fmt.Fprintf(w,"405 - %s Method not allowed",r.Method)
		return
	}

	query := r.URL.Query()
	opt := options.Find()

	offset, present := query["offset"]
	var p1_offset int64 = 0
	var p2_limit int64 = 0
	p3_sortOps := bson.D{}

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

		p1_offset = oset
	}

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
	
	p2_limit = lmt

	err := utils.AddFilters(query, opt)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w,"Error: Bad Request")
		return
	}

	filters, present := query["filters"]

	if present {
		filter_list, err := utils.CheckFilter(filters[0])
		
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w,"Error: Bad Request")
			return
		}


		for _, filter := range filter_list {
			p3_sortOps = append(p3_sortOps, bson.E{Key: filter , Value: 1})
		}

		opt.SetSort(p3_sortOps)
	}

	
	// pipeline stages
	pipeline := bson.A {
		bson.D{{Key: "$skip", Value: p1_offset}},
		bson.D{{Key: "$limit", Value: p2_limit}},
	}

	if len(p3_sortOps)!=0 {
		pipeline = append(pipeline,bson.D{{Key: "$sort", Value: p3_sortOps}} ) 
	}


	// aggregate
	cursor, err := collection.Aggregate(context.TODO(),pipeline)

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
	if err = cursor.All(context.TODO(), &result); err != nil {
  		log.Fatal(err)
	}

	// return json object through HTTP response
	json.NewEncoder(w).Encode(result) 
}

// search anime based on title
func SearchAnime(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		fmt.Fprintf(w,"405 - %s Method not allowed",r.Method)
		return
	}
	
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