FROM golang:1.19

WORKDIR /app

# build the main app
COPY ./src /src
WORKDIR /src
RUN go get \
&&  go build -o /app/main -a -v /src/main.go

# start the app
WORKDIR /app
COPY ./entrypoint.sh /app
ENTRYPOINT [ "/app/entrypoint.sh" ]
