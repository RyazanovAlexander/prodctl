FROM busybox:1.35.0 AS bundle

WORKDIR /bundle/
COPY . .
RUN rm ./Dockerfile

# -----------------------------------------------

FROM gcr.io/distroless/static:nonroot
COPY --from=bundle /bundle/ /bundle/
USER nonroot:nonroot