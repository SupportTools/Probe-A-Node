FROM golang:1.20 AS build_base
ARG GIT_COMMIT=unknown
ARG GIT_BRANCH=unknown
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=$GIT_COMMIT" -o main main.go
RUN chmod +x main

FROM scratch
COPY --from=build_base /src/main /usr/local/bin/Probe-A-Node
CMD ["/usr/local/bin/Probe-A-Node"]