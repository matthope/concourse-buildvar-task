FROM --platform=$BUILDPLATFORM golang:1.21.4-alpine3.17 AS base
WORKDIR /opt/resource

FROM base AS build
WORKDIR /src

ENV GOOS=$TARGETOS GOARCH=$TARGETARCH

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -ldflags="-w -s" -o /task .

FROM scratch
LABEL org.opencontainers.image.source https://github.com/matthope/concourse-buildvar-task

COPY --from=build /task /task

CMD [ "/task" ]
