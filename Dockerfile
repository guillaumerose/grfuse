FROM golang:alpine
RUN apk add --no-cache fuse
COPY . /go/src/github.com/LK4D4/grfuse
RUN go build github.com/LK4D4/grfuse/example/helloclient
RUN mkdir /go/foo
CMD ./helloclient