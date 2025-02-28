FROM golang:1.24

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download && go mod verify  && go mod tidy

COPY . .

RUN go build -v -o server

CMD ["/app/server"]