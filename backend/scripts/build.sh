rm -rf bin/
mkdir -p bin
docker run -v ${PWD}:/go/src/github.com/foundcenter/moas/backend/ --rm golang:1.8.0 /bin/sh -c 'cd /go/src/github.com/foundcenter/moas/backend/ && go build -o bin/moas .'
docker build -t moas:latest .