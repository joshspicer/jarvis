
FROM golang:1.18

ARG JARVIS_BUILD_COMMIT="0000"
ARG JARVIS_BUILD_VERSION="dev"

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build \ 
    --ldflags "-X 'main.commit=$JARVIS_BUILD_COMMIT' -X 'main.version=$JARVIS_BUILD_VERSION'" \
     -v -o /usr/local/bin/app ./...

CMD ["app"]