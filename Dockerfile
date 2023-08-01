FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o application cmd/main.go

EXPOSE 3030

CMD [ "/app/application" ]

