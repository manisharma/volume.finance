FROM golang:alpine as builder
ARG PORT=8080
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod download 
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/cmd/server .
EXPOSE ${PORT}
ENTRYPOINT [ "./server" ]