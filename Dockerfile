FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/kbot .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S kbot \
    && adduser -S -G kbot kbot

COPY --from=builder /out/kbot /usr/local/bin/kbot

USER kbot

ENTRYPOINT ["/usr/local/bin/kbot", "start"]