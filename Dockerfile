FROM golang:1.19.2-alpine as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./...

FROM alpine

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
WORKDIR /home/nonroot

RUN apk add -U --no-cache ca-certificates openjdk11-jre

COPY --from=builder /usr/local/bin/aws-kcl /usr/local/bin/aws-kcl
COPY --from=builder --chown=nonroot:nonroot /usr/local/bin/sample ./sample/sample
COPY --chown=nonroot:nonroot pom.xml sample/development.properties ./

USER nonroot:nonroot

CMD aws-kcl -java /usr/bin/java -properties development.properties



