FROM golang:1.19

# set workdir
WORKDIR /app

# build the main app
COPY ./src /src
RUN go build -a -v /src/main.go

# start the app
COPY ./entrypoint.sh /app
ENTRYPOINT [ "/app/entrypoint.sh" ]
