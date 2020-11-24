# Owner's Manual

## Technical Notes

### Preserving Key Order of JSON

On the back-end, when using the provided JSON decoder to unmarshal a json string to 
a map[string]interface{} data type, the order of the keys will be lost when iterating
over all the keys in the resulting map.

To fix this, a package called "go-ordered-json" is used which retains the order of keys.

Here is some more info about that package
* Code: https://gitlab.com/c0b/go-ordered-json
* Article: https://medium.com/@ty0h/preserving-json-object-keys-order-in-javascript-python-and-go-language-170eaae0de03
* Docs: https://godoc.org/gitlab.com/c0b/go-ordered-json

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
