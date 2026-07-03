FROM golang:1.25 AS builder
#FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

# --- for alpine:
#RUN apk --no-cache add bash git make gcc gettext musl-dev
COPY ["go.mod","go.sum","./"]
RUN go mod download
COPY ./ ./
RUN go build -o main cmd/sqlx/main_sqlx.go

# --- https://github.com/GoogleCloudPlatform/golang-samples/blob/main/appengine/go11x/tasks/handle_task/Dockerfile
# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM ubuntu:24.04
#FROM golang
#FROM alpine
COPY --from=builder /usr/local/src/main /apps/main
#COPY db/todolist.sqlite /apps/todolist.sqlite
#COPY front-react/static /apps/static
#COPY ./dist /apps/dist
#COPY front-react/index.html /apps/index.html
#COPY app/docs	/apps/docs

#VOLUME /apps/web/upload

WORKDIR /apps
ENV PATH="/apps:${PATH}"
CMD ["tail", "-f", "/dev/null"]
ENTRYPOINT ["main"]