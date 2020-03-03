#!/bin/sh
cmd="go run . -env dev -workDir .\\\\sample $@ sample.yaml"
sh -c "$cmd"

