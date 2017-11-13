#!/bin/bash

TAG=$(git rev-parse --short HEAD)
docker push andrewbackes/tourney:${TAG}