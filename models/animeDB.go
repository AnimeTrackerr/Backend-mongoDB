package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnimeSeason struct {
	Season	string 	`bson:"season"`
	Year	int32 	`bson:"year"`
}

type Anime struct {
	ID          primitive.ObjectID  `bson:"_id"`
	Sources     []string 			`bson:"sources"`
	Title       string   			`bson:"title"`
	Type        string   			`bson:"type"`
	Episodes    int32    			`bson:"episodes"`
	Status      string   			`bson:"status"`
	AnimeSeason AnimeSeason 		`bson:"animeSeason"`
	Picture		string 				`bson:"picture"`
	Thumbnail 	string 				`bson:"thumbnail"`
	Synonyms 	string 				`bson:"synonyms"`
	Relations   []string 			`bson:"relations"`
	Tags		[]string 			`bson:"tags"`
	MalID		string 				`bson:"malID"`
	MalScore  	float64				`bson:"malScore,truncate"`
	StartDate	bson.M  			`bson:"startDate"`
	EndDate		bson.M				`bson:"endDate"`
	Nsfw		string 				`bson:"nsfw"`
	Rating		string 				`bson:"rating"`
	Synopsis	string 				`bson:"synopsis"`

}
