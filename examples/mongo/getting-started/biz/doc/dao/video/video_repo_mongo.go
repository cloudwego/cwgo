// Code generated by cwgo (v0.0.1). DO NOT EDIT.

package video

import (
	"context"
	"github.com/cloudwego/cwgo/examples/mongo/getting-started/biz/doc/model/video"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewVideoRepository(collection *mongo.Collection) VideoRepository {
	return &VideoRepositoryMongo{
		collection: collection,
	}
}

type VideoRepositoryMongo struct {
	collection *mongo.Collection
}

func (r *VideoRepositoryMongo) InsertVideo(ctx context.Context, video *video.Video) (interface{}, error) {
	result, err := r.collection.InsertOne(ctx, video)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}