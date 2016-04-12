FROM golang:1.6
ADD . /go/src/github.com/bbokorney/pipeline
WORKDIR /go/src/github.com/bbokorney/pipeline
ENV GOPATH=/go:/go/src/github.com/bbokorney/pipeline/Godeps/_workspace
RUN go build -o /app
CMD ["/app"]
