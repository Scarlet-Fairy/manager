package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Scarlet-Fairy/manager/pb"
	"github.com/Scarlet-Fairy/manager/pkg/endpoint"
	"github.com/Scarlet-Fairy/manager/pkg/logger"
	amqpMessage "github.com/Scarlet-Fairy/manager/pkg/message/amqp"
	mongoRepository "github.com/Scarlet-Fairy/manager/pkg/repository/mongo"
	grpcScheduler "github.com/Scarlet-Fairy/manager/pkg/scheduler/grpc"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	grpcTransport "github.com/Scarlet-Fairy/manager/pkg/transport/grpc"
	"github.com/go-kit/kit/log/level"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/run"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

var (
	grpcAddr      = flag.String("grpc-url", ":8081", "gRPC server listen address")
	schedulerUrl  = flag.String("scheduler-url", "localhost:8082", "url of the scheduler gRPC server")
	amqpUrl       = flag.String("amqp-url", "amqp://demo:demo@localhost:5672", "url of rabbitmq")
	mongoUrl      = flag.String("mongo-url", "mongodb://localhost:27017", "url of mongodb")
	mongoDatabase = flag.String("mongo-db", "manager", "mongodb manager where store state data")
)

var (
	loggers                   = logger.NewLogger()
	infoLogger                = loggers.InfoLogger
	warnLogger                = loggers.WarnLogger
	errorLogger               = loggers.ErrorLogger
	transportLayerLogger      = loggers.TransportLayerLogger
	endpointLayerLogger       = loggers.EndpointLayerLogger
	coreLayerLogger           = loggers.CoreLayerLogger
	serviceComponentLogger    = loggers.ServiceComponentLogger
	messageComponentLogger    = loggers.MessageComponentLogger
	repositoryComponentLogger = loggers.RepositoryComponentLogger
	schedulerComponentLogger  = loggers.SchedulerComponentLogger
)

var ctx = context.Background()

func main() {
	flag.Parse()

	schedulerClient, err := newSchedulerClient(*schedulerUrl)
	if err != nil {
		errorLogger.Log(
			"scheduler-url", *schedulerUrl,
			"during", "init",
			"msg", "scheduler gRPC client init failed",
			"err", err,
		)
		os.Exit(1)
	}

	rabbitMqClient, err := newRabbitMQClient(*amqpUrl)
	if err != nil {
		errorLogger.Log(
			"amqp-url", *amqpUrl,
			"during", "init",
			"msg", "amqp client init failed",
			"err", err,
		)
		os.Exit(1)
	}
	defer func() {
		if err := rabbitMqClient.Close(); err != nil {
			panic(err)
		}
	}()
	rabbitMqChannel, err := rabbitMqClient.Channel()
	if err != nil {
		errorLogger.Log(
			"amqp-url", *amqpUrl,
			"during", "init",
			"msg", "amqp channel creation failed",
			"err", err,
		)
		os.Exit(1)
	}
	defer func() {
		if err := rabbitMqChannel.Close(); err != nil {
			panic(err)
		}
	}()

	mongoDbClient, err := newMongoDbClient(ctx, *mongoUrl)
	if err != nil {
		errorLogger.Log(
			"mongo-url", *mongoUrl,
			"during", "init",
			"msg", "mongodb client init failed",
			"err", err,
		)
		os.Exit(1)
	}
	defer func() {
		if err := mongoDbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	deployCollection := mongoDbClient.Database(*mongoDatabase).Collection("deploy")

	amqpMessageInstance := amqpMessage.New(rabbitMqChannel)
	if err := amqpMessageInstance.Init(); err != nil {
		level.Error(messageComponentLogger).Log(
			"during", "init",
			"msg", "Init failed",
			"err", err,
		)
		os.Exit(1)
	}

	mongoRepositoryInstance := mongoRepository.New(deployCollection)

	grpcSchedulerInstance := grpcScheduler.New(schedulerClient)

	svc := service.NewService(mongoRepositoryInstance, amqpMessageInstance, grpcSchedulerInstance)
	endpoints := endpoint.NewEndpoint(svc, endpointLayerLogger)
	grpcServer := grpcTransport.NewGRPCServer(endpoints, transportLayerLogger)

	var g run.Group
	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			level.Error(transportLayerLogger).Log(
				"during", "init-listen",
				"msg", fmt.Sprintf("failed to listen on %s", *grpcAddr),
				"err", err,
			)
			os.Exit(1)
		}

		g.Add(func() error {
			transportLayerLogger.Log(
				"addr", *grpcAddr,
			)

			baseServer := grpc.NewServer(
				grpc.UnaryInterceptor(kitgrpc.Interceptor),
			)
			pb.RegisterManagerServer(baseServer, grpcServer)

			return baseServer.Serve(grpcListener)
		}, func(err error) {
			if err = grpcListener.Close(); err != nil {
				panic(err)
			}
		})
	}

	infoLogger.Log("exit", g.Run())
}

func newSchedulerClient(url string) (pb.SchedulerClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewSchedulerClient(conn)
	return client, err
}

func newRabbitMQClient(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func newMongoDbClient(ctx context.Context, url string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := client.Connect(timeoutCtx); err != nil {
		return nil, err
	}

	if err := client.Ping(timeoutCtx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
