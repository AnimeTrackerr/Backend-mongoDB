package utils

import (
	"errors"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var FilterOps map[string]bool = map[string]bool {
	"malScore" : true,
}

// check query params
func CheckQueryParam(pararms string, expected_type, err error) {
}


// check filter options
func CheckFilter(filters string) ([]string, error) {
	var filter_list []string = strings.Split(filters, ",")

	for _, filter := range filter_list {
		_, ok := FilterOps[filter] 
		
		if  !ok {
			return filter_list, errors.New("Key Error: " + filter)
		} 
	}
	
	return filter_list, nil
}

// add filter options
func AddFilters(query url.Values, findops *options.FindOptions) error {
	filters, present := query["filters"]

	if present {
		filter_list, err := CheckFilter(filters[0])
		
		if err != nil {
			return err
		}

		sortOps := bson.D{}

		for _, filter := range filter_list {
			sortOps = append(sortOps, bson.E{Key: filter , Value: 1})
		}

		findops.SetSort(sortOps)
	}

	return nil
}
