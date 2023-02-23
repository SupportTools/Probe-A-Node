FROM golang:1.20 as build

ARG GIT_COMMIT
WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o probe-a-node main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/app/probe-a-node /app/
ENTRYPOINT /app
EXPOSE 9876
