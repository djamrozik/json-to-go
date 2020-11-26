test-server:
	cd server && go test -v ./lib

build-server: test-server add-client-to-server
	cd server && go build main.go

run-server: test-server add-client-to-server
	cd server && modd

build-client:
	cd client && yarn build

run-client:
	cd client && yarn start

add-client-to-server: build-client
	cp -R client/build server/
