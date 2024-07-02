#!/usr/bin/env bash

awslocal dynamodb create-table \
  --endpoint-url=http://localhost:4566 \
  --table-name singletable \
  --region=us-east-1 \
  --attribute-definitions AttributeName=PK,AttributeType=S AttributeName=SK,AttributeType=S \
  --key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --table-class STANDARD