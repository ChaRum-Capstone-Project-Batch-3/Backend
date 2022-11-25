FROM golang:alpine

WORKDIR /app

COPY . /app

RUN go mod download
RUN go build -o dist

EXPOSE 443

ENTRYPOINT [ "./dist" ]