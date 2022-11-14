package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ClientInfo struct {
	Client *mongo.Client
	Ctx context.Context
}