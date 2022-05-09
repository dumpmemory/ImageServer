all: build

build:
	export GIN_MODE=release; \
	go build -o image_server -ldflags "-X app.Version=`git describe --tags --abbrev=0 --dirty=-dev` -X app.Build=`git rev-parse --short HEAD`" main.go

run:
	go run main.go

clean:
	rm image_server