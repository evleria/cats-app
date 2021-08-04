FROM golang:1.16-alpine AS build

WORKDIR /
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go /
COPY /internal /internal
RUN go build -o server ./main.go

FROM alpine

COPY --from=build /server /server

CMD ["./server"]