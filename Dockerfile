FROM golang:1.26rc1-alpine3.23  AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o server ./cmd/comp-math-6/comp-math-6.go

FROM alpine:3.23

COPY --from=builder /app/server /app/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]
