FROM golang:1.18 AS builder
WORKDIR /go/src/github.com/ThelonKarrde/Kube-Image-Guardian
COPY go.mod go.sum ./
RUN go get -d -v -u all
COPY . /go/src/github.com/ThelonKarrde/Kube-Image-Guardian
RUN go build -o ./kube-image-guardian

FROM alpine:3.15.4
WORKDIR /app
RUN adduser -h /app -D web
COPY --from=builder /go/src/github.com/ThelonKarrde/Kube-Image-Guardian/kube-image-guardian /app/
USER web
ENTRYPOINT ["/app/kube-image-guardian"]