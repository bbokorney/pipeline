FROM golang:1.6

ADD pipeline /pipeline
CMD ["/pipeline"]
