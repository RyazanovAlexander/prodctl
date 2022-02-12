ARG BUILD_IMAGE_TAG=latest
FROM golang:${BUILD_IMAGE_TAG} AS builder

WORKDIR /src/
COPY . .

ARG GOOS=linux
ARG GOARCH=amd64
ARG LDFLAGS="-w -s"
RUN go get -d -v
RUN GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o /bin/prodctl

# -----------------------------------------------

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /bin/prodctl /bin/prodctl
USER nonroot:nonroot
ENTRYPOINT ["/bin/prodctl"]