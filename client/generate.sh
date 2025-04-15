#!/bin/bash
go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.yaml openapi.json
