unit-test:
	go test ./... -v

clean:
	rm -rf bin

build:clean
	mkdir bin
	go build -o bin/main

run:
	API_KEY="f9adf6ce-9aa0-4071-a768-ac0525f2a966" \
	MONGO_USERNAME=root \
	MONGO_PASSWORD=secretpass \
	MONGO_URI="mongodb://localhost:27017" \
	go run main.go setup.go
