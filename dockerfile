### STEP 1: Builder ###

FROM golang:1.26-alpine AS builder


WORKDIR /build

#Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

#Copy rest of the code
COPY . .

#Compile the binary
RUN CGO_ENABLED=0 go build -o /build/api ./cmd/api

### STEP 2: Runner ###

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/api .

EXPOSE 3000

CMD ["./api"]