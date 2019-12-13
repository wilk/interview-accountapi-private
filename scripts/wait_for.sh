#! /usr/bin/env bash

# sync with the accountapi service using netcat
echo "Attempting to connect to accountapi"
until $(nc -zv accountapi 8080); do
    printf '.'
    sleep 5
done
echo "Was able to connect to accountapi. Proceeding with tests"

go test
