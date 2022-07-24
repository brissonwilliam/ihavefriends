# ********************************
# DOCKER BUILD ARGS
# BuiltKit must be enabled.
#   DOCKER_BUILDKIT=1 docker build .
#   or modify /etc/docker/daemon.json to have { "features": { "buildkit": true } }
#   or JetBrains IDE has an option under Edit Configuration -> Build option, click the modify options link label -> Enable BuiltKit
#
# DOCKER RUN ARGS
# --env-file=path_to_env
#   for injecting the env vars. Or else, it will default to dev env vars.
# --network host
#   FOR LINUX ONLY to allow the docker to connect to the local machine ports (mainly for database)
#   For windows, set the env var DB_HOST to host.docker.internal in dev.env or through -e


# Ports must be exposed for the web server and database with the -p arg
# ********************************

# ********************************
# STEP 1
# Build app and dependencies
# ********************************
FROM golang:1.17 as build
ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org
WORKDIR /code

# download and build dependencies. Do this before copying other source files to speed up the build
COPY go.mod go.sum /code/
RUN go mod download

# copy source code and build app with cache dependencies
COPY . /code
RUN --mount=type=cache,target=/root/.cache/go-build go build -o backend

# ********************************
# STEP 2
# Copy build outputs and run in alpine
# ********************************
FROM alpine:3.15.4

RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /code/backend /backend
COPY --from=build /code/dev.env /dev.env

ENTRYPOINT [ "/backend" ]