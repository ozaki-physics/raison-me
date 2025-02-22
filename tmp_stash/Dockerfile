# Google App Engine の最新
FROM golang:1.19

ENV REPOSITORY=github.com/ozaki-physics/raison-me
WORKDIR /go/src/$REPOSITORY
COPY . .
RUN go mod download
# RUN go build -o /go/bin main.go
# CMD ["/go/bin/main"]
