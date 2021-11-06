FROM golang:1.17-alpine
WORKDIR /go/src/github.com/ThelonKarrde/Kube-Image-Guardian
COPY go.mod go.sum ./
RUN go get -d -v -u all
COPY . /go/src/github.com/ThelonKarrde/Kube-Image-Guardian
RUN go build


FROM alpine:3.13.6
WORKDIR /app
RUN adduser -h /app -D web
COPY --from=0 /go/src/github.com/ThelonKarrde/Kube-Image-Guardian/Kube-Image-Guardian /app/
USER web
ENTRYPOINT ["./Kube-Image-Guardian"]
