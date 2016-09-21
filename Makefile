default: build

lint:
	go fmt 
	go vet 
	gometalinter --deadline=15s ./...

build:
	go fmt 
	go vet 
	go build 
test: build
	go test -v ./...
   zip: build	
	zip invstick-linux-$(shell arch).zip invstick