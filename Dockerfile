FROM golang:1.16-alpine AS build

WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go /
COPY /internal /internal
RUN go build -o server ./main.go

FROM alpine

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

COPY --from=build /server /server

CMD ["./server"]