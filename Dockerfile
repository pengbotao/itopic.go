FROM golang:1.14 AS build

COPY . /go/src/github.com/pengbotao/itopic.go/
RUN cd /go/src/github.com/pengbotao/itopic.go/ \
&& CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o itopic


FROM alpine
COPY --from=build /go/src/github.com/pengbotao/itopic.go/ /www/itopic.go/
EXPOSE 8001
WORKDIR /www/itopic.go
CMD [ "-host", "0.0.0.0:8001" ]
ENTRYPOINT [ "/www/itopic.go/itopic" ]