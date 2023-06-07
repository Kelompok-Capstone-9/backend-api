FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app

RUN go build -o main .

ENV PORT=8000

EXPOSE 8000

CMD [ "./main" ]