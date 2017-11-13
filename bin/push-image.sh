#!/bin/bash -x

TAG=$(git rev-parse --short HEAD)
docker push andrewbackes/tourney:${TAG}