FROM golang:1.9
 
WORKDIR /go/src/github.com/kernkw/hhapp

ENV HHAPP_DB_HOST=localhost \
    HHAPP_DB_NAME=hhapp \
    HHAPP_DB_USER=root \
    HHAPP_DB_PASSWORD= 

### BUILD ###

# install vendored files (very rarely change) early for quick
# build speed and so that the layer can be cached.
COPY vendor ./vendor

# copy the rest of the repo root into the container workdir
COPY . .

# pre-build runtime and test dependencies. This makes any testing fast
# in the presence of slow-to-compile packages (e.g. cgo)
RUN go install ./... && go install -tags test -race ./...

CMD ["./bin/start"]

EXPOSE 8080
