#!/bin/sh
go run . -name "My first deployment" -dir ./.build -hosts ./.hosts -vars ./vars.yaml  sample.yaml
