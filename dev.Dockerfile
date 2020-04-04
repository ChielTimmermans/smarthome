FROM golang:1.13

# Set the Current Working Directory inside the container
WORKDIR /smarthome-home

COPY . .

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -log-prefix=false -exclude-dir=.git -build="go build -o smarthome-home ./cmd/api" -command="./smarthome-home -logtostderr"
