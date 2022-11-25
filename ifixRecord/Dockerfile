FROM golang:alpine As builder
RUN mkdir /recordbuild
ADD . /recordbuild
WORKDIR /recordbuild
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/natefinch/lumberjack
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /recordbuild/ /app
CMD ["/app/main"]
