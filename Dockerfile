ARG BUILD_IMAGE_TAG=latest
FROM golang:${BUILD_IMAGE_TAG} AS builder

WORKDIR /src/
COPY . .

ARG GOOS=linux
ARG GOARCH=amd64
ARG LDFLAGS="-w -s"
RUN go get -d -v
RUN GOOS=$GOOS GOARCH=$GOARCH GO111MODULE=on go build -ldflags "$LDFLAGS" -o /bin/command-executor

# -----------------------------------------------

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /bin/command-executor /bin/command-executor
USER nonroot:nonroot
ENTRYPOINT ["/bin/command-executor"]