FROM golang:1.8 as builder
WORKDIR $GOPATH/src/github.com/brightsparc/ndc-golang-talk
RUN go get
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder $GOPATH/src/github.com/brightsparc/ndc-golang-talk .
EXPOSE 3001
CMD ["./app"]
