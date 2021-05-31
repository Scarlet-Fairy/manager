package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Env struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

type Workload struct {
	JobId   string `bson:"job_id"`
	JobName string `bson:"job_name"`
	Envs    []*Env `bson:"envs"`
	Url     string `bson:"url"`
}

type BuildStep struct {
	Step  int    `bson:"step"`
	Error string `bson:"error"`
}

type Build struct {
	JobId     string       `bson:"job_id"`
	JobName   string       `bson:"job_name"`
	ImageName string       `bson:"image_name"`
	Status    int          `bson:"status"`
	Steps     []*BuildStep `bson:"steps"`
}

type Deploy struct {
	Id       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	GitRepo  string             `bson:"git_repo"`
	Build    *Build             `bson:"build"`
	Workload *Workload          `bson:"workload"`
}
