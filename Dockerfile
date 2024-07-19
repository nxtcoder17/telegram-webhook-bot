FROM gcr.io/distroless/static
# FROM alpine:3.18
ARG BINARY
COPY ${BINARY} /tg-webhook-bot
ENTRYPOINT [ "/tg-webhook-bot" ]
