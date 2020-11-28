test-server:
	cd server && go test -v ./lib

build-server: test-server add-client-to-server
	cd server && go build main.go

run-server: test-server add-client-to-server
	cd server && modd

build-client:
	cd client && yarn install && yarn build

run-client:
	cd client && yarn start

add-client-to-server: build-client
	cp -R client/build server/

build-image:
	docker build . -t json-to-golang:latest

tag-image:
	docker tag json-to-golang:latest gcr.io/json-to-golang/json-to-golang:latest

upload-image: build-image tag-image
	docker push gcr.io/json-to-golang/json-to-golang:latest
