FROM golang:1.18-alpine AS BUILDER

LABEL stage backendbuilder

ARG GOOS
ENV GOOS=${GOOS:-linux}
ARG GOARCH
ENV GOARCH=${GOARCH:-amd64}
ARG GOPROXY
ENV GOPROXY=${GOPROXY:-https://goproxy.io}

WORKDIR /app

COPY . /app
RUN go mod download
RUN go build -o homework main.go

FROM alpine:3

LABEL maintainer "Fredy Yu"
LABEL org.label-schema.name "Homework"
LABEL org.label-schema.vendor "Fredy Yu"
LABEL org.label-schema.schema-version "0.0.1"

WORKDIR /app

RUN set -eux; \
    apk update \
    && apk upgrade \
    && apk add --no-cache \
    dumb-init \
    curl

COPY --from=BUILDER /app/homework /app/

ENTRYPOINT [ "dumb-init", "--" ]

CMD [ "/app/homework", "start" ]

EXPOSE 80
