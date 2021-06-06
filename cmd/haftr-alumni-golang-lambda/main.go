package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/BenBraunstein/haftr-alumni-golang/internal/app"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type awsEventHandlerFunc func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func main() {
	h := getAWSLambdaEventHandler()
	lambda.Start(h)
}

func getAWSLambdaEventHandler() awsEventHandlerFunc {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		mongoURI := os.Getenv("MONGO_URI")
		dbName := os.Getenv("DB_NAME")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(errors.Wrap(err, "main - cannot connect to mongo"))
		}
		defer client.Disconnect(ctx)
		db := client.Database(dbName)

		a := app.New(db)
		lambdaProxyAdapter := handlerfunc.New(a.Handler())
		return lambdaProxyAdapter.ProxyWithContext(ctx, req)
	}
}
