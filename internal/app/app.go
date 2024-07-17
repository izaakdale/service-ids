package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/ptr"
	"github.com/go-redis/redis"
	dsdynamo "github.com/izaakdale/service-ids/internal/datastore/dynamo"
	dsredis "github.com/izaakdale/service-ids/internal/datastore/redis"
	"github.com/izaakdale/service-ids/internal/router"
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Host          string `envconfig:"HOST"`
	Port          int    `envconfig:"PORT" default:"80"`
	AWSRegion     string `envconfig:"AWS_REGION" default:"us-east-1"`
	AWSAcessKeyID string `envconfig:"AWS_ACCESS_KEY_ID"`
	AWSSecretKey  string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	AWSEndpoint   string `envconfig:"AWS_ENDPOINT"`
	TableName     string `envconfig:"TABLE_NAME"`
	RedisEndpoint string `envconfig:"REDIS_ENDPOINT"`
	UseDynamo     bool   `envconfig:"USE_DYNAMO" default:"false"`
	UseRedis      bool   `envconfig:"USE_REDIS" default:"false"`
}

func Run() error {
	var spec specification
	envconfig.MustProcess("", &spec)
	log.Printf("%+v\n", spec)

	errCh := make(chan error, 1)
	if spec.UseDynamo {
		cfg, err := config.LoadDefaultConfig(context.Background(), func(lo *config.LoadOptions) error {
			if spec.AWSRegion != "" {
				lo.Region = spec.AWSRegion
			}

			if spec.AWSAcessKeyID != "" && spec.AWSSecretKey != "" {
				lo.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{
						AccessKeyID: spec.AWSAcessKeyID, SecretAccessKey: spec.AWSSecretKey,
					}, nil
				})
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
		cli := dsdynamo.New(svc, spec.TableName)
		mux := router.New(cli)

		log.Println("starting dynamo server...")
		go func() {
			errCh <- http.ListenAndServe(fmt.Sprintf("%s:%d", spec.Host, spec.Port), mux)
		}()
	}

	if spec.UseRedis {
		opt, err := redis.ParseURL(spec.RedisEndpoint)
		if err != nil {
			return err
		}
		redCli := redis.NewClient(opt)

		for i := 0; i < 100; i++ {
			redCli.HSet(fmt.Sprintf("bing-%d", i), "harvester", fmt.Sprintf("harvester-%d", i))
			redCli.HSet(fmt.Sprintf("bing-%d", i), "some", fmt.Sprintf("other-%d", i))
			bytes, _ := json.Marshal(struct {
				One string `json:"one"`
			}{
				One: fmt.Sprintf("bing-%d", i),
			})
			b64 := base64.StdEncoding.EncodeToString(bytes)
			redCli.HSet(fmt.Sprintf("bing-%d", i), "payload", b64)
		}
		for i := 0; i < 100; i++ {
			redCli.HSet(fmt.Sprintf("harvester-%d", i), "bing", fmt.Sprintf("bing-%d", i))
		}

		cli2 := dsredis.New(redCli, spec.TableName)
		mux2 := router.New(cli2)

		log.Println("starting redis server...")
		go func() {
			errCh <- http.ListenAndServe(fmt.Sprintf("%s:%d", spec.Host, spec.Port+1), mux2)
		}()
	}

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
