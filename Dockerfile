FROM golang:1.22

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

RUN go mod tidy

COPY . .

EXPOSE 9000

CMD [ "go", "run", "main.go" ]