package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/ptr"
	"github.com/izaakdale/service-ids/internal/datastore"
	"github.com/izaakdale/service-ids/internal/router"
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Host        string `envconfig:"HOST"`
	Port        int    `envconfig:"PORT" default:"80"`
	AWSRegion   string `envconfig:"AWS_REGION" default:"us-east-1"`
	AWSEndpoint string `envconfig:"AWS_ENDPOINT"`
	TableName   string `envconfig:"TABLE_NAME" required:"true"`
}

func Run() error {
	var spec specification
	envconfig.MustProcess("", &spec)
	log.Println(spec)

	cfg, err := config.LoadDefaultConfig(context.Background(), func(lo *config.LoadOptions) error {
		if spec.AWSRegion != "" {
			lo.Region = spec.AWSRegion
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	svc := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if spec.AWSEndpoint != "" {
			o.BaseEndpoint = ptr.String(spec.AWSEndpoint)
		}
	})
	cli := datastore.New(svc, spec.TableName)
	mux := router.New(cli)
	log.Println("starting server...")

	errCh := make(chan error, 1)
	go func() {
		errCh <- http.ListenAndServe(fmt.Sprintf("%s:%d", spec.Host, spec.Port), mux)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	select {
	case <-signalCh:
		log.Println("shutting down...")
	case err := <-errCh:
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}
