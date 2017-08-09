FROM golang:1.8 as builder
WORKDIR /go/src/github.com/brightsparc/ndc-golang-talk/ndc_segment
COPY config.json .
COPY main.go .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ndc_segment .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/brightsparc/ndc-golang-talk/ndc_segment .
EXPOSE 3001
CMD ["./ndc_segment"]
