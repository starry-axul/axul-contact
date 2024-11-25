FROM golang:1.21

WORKDIR $GOPATH/bin

COPY main .

EXPOSE 8081

CMD ["main"]