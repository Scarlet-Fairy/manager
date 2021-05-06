package mongo

import (
	"context"
	middlewares "github.com/Scarlet-Fairy/manager/pkg/repository"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection, logger log.Logger) service.Repository {
	var instance service.Repository
	instance = &mongoRepository{
		collection: collection,
	}
	instance = middlewares.LoggingMiddleware(logger)(instance)

	return instance
}

func (m *mongoRepository) CreateDeploy(ctx context.Context, deploy *service.Deploy) (string, error) {
	res, err := m.collection.InsertOne(ctx, businessToData(deploy))
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).String(), nil
}

func (m *mongoRepository) GetDeploy(ctx context.Context, id string) (*service.Deploy, error) {
	res := m.collection.FindOne(
		ctx,
		bson.M{
			"_id": bson.M{
				"$eq": id,
			},
		},
	)
	if err := res.Err(); err != nil {
		return nil, err
	}

	deploy := &Deploy{}
	if err := res.Decode(deploy); err != nil {
		return nil, err
	}

	return dataToBusiness(deploy), nil
}

func (m *mongoRepository) ListDeploy(ctx context.Context) ([]*service.Deploy, error) {
	var deploys []*service.Deploy

	cur, err := m.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var dataDeploy *Deploy

		if err := cur.Decode(dataDeploy); err != nil {
			return nil, err
		}

		deploys = append(deploys, dataToBusiness(dataDeploy))
	}

	return deploys, nil
}

func (m *mongoRepository) UpdateDeploy(ctx context.Context, deploy *service.Deploy) error {
	dataDeploy := businessToData(deploy)

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": bson.M{
				"$eq": dataDeploy.Id,
			}},
		bson.M{
			"$set": dataDeploy,
		},
	)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (m *mongoRepository) DeleteDeploy(ctx context.Context, id string) error {
	dataId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := m.collection.DeleteOne(ctx, bson.M{
		"_id": bson.M{
			"$eq": dataId,
		},
	})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("0 Doc has been deleted")
	}

	return nil
}

func (m *mongoRepository) InitBuild(ctx context.Context, id string, jobName string, jobId string, imageName string) error {
	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": bson.M{
				"$eq": id,
			},
		},
		bson.M{
			"$set": Deploy{
				Build: &Build{
					JobId:     jobId,
					JobName:   jobName,
					ImageName: imageName,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (m *mongoRepository) InitWorkload(ctx context.Context, id string, jobName string, jobId string, envs map[string]string) error {
	dataEnv := mapEnvToArrEnv(envs)

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": bson.M{
				"$eq": id,
			},
		},
		bson.M{
			"$set": Deploy{
				Workload: &Workload{
					JobId:   jobId,
					JobName: jobName,
					Envs:    dataEnv,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (m *mongoRepository) SetBuildStatus(ctx context.Context, id string, status service.Status) error {
	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": bson.M{
				"$eq": id,
			},
		},
		bson.M{
			"$set": Deploy{
				Build: &Build{
					Status: int(status),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}

func (m *mongoRepository) RecordBuildStep(ctx context.Context, id string, buildStep service.BuildStep) error {
	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": bson.M{
				"$eq": id,
			},
		},
		bson.M{
			"$push": bson.M{
				"build": bson.M{
					"steps": BuildStep{
						Step:  int(buildStep.Step),
						Error: buildStep.Error,
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("document not found")
	}

	return nil
}
