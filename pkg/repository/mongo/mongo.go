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

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *mongoRepository) GetDeploy(ctx context.Context, id string) (*service.Deploy, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := m.collection.FindOne(
		ctx,
		bson.M{
			"_id": objectId,
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

	cur, err := m.collection.Find(ctx, bson.M{})
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
	objectId, err := primitive.ObjectIDFromHex(deploy.Id)
	if err != nil {
		return err
	}

	dataDeploy := businessToData(deploy)

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": objectId,
		},
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := m.collection.DeleteOne(ctx, bson.M{
		"_id": objectId,
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": objectId,
		},
		bson.M{
			"$set": bson.M{
				"build.job_id":     jobId,
				"build.job_name":   jobName,
				"build.image_name": imageName,
				"build.steps":      []bson.M{},
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	dataEnv := mapEnvToArrEnv(envs)

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": objectId,
		},
		bson.M{
			"$set": bson.M{
				"workload.job_id":   jobId,
				"workload.job_name": jobName,
				"workload.envs":     dataEnv,
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": objectId,
		},
		bson.M{
			"$set": bson.M{
				"build.status": int(status),
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := m.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": objectId,
		},
		bson.M{
			"$push": bson.M{
				"build.steps": bson.M{
					"step":  int(buildStep.Step),
					"error": buildStep.Error,
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
