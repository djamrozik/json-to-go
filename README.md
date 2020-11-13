# Owner's Manual

## Building and Deploying to Container Register

Prerequisites for building/deploying
* Docker is installed
* GCP SDK is installed
* Configure Docker to use GCP with `gcloud auth configure-docker`
* Docker daemon is running

To build run the following command from the main directory:
* `docker build . -t json-to-golang:latest`

Deploying
* Tag the image `docker tag json-to-golang:latest gcr.io/json-to-golang/json-to-golang:latest`
* Push to container registry `docker push gcr.io/json-to-golang/json-to-golang:latest`
