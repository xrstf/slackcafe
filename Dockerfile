FROM golang:1.12-alpine as builder

RUN apk add --update make git
WORKDIR /go/src/github.com/xrstf/slackcafe
COPY . .
RUN make build

FROM alpine:3.10

RUN apk --no-cache add ca-certificates tzdata tesseract-ocr tesseract-ocr-data-deu

ENTRYPOINT ["./slackcafe"]

WORKDIR /app
COPY --from=builder /go/src/github.com/xrstf/slackcafe/slackcafe .
