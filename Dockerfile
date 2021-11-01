FROM golang:1.17-alpine

COPY . /go/src/github.com/ThelonKarrde/Kube-Image-Guardian
WORKDIR /go/src/github.com/ThelonKarrde/Kube-Image-Guardian
RUN go get
RUN go build


FROM alpine
WORKDIR /app
RUN adduser -h /app -D web
COPY --from=0 /go/src/github.com/ThelonKarrde/Kube-Image-Guardian/Kube-Image-Guardian /app/

USER web
ENTRYPOINT ["./Kube-Image-Guardian"]
EXPOSE 1224
