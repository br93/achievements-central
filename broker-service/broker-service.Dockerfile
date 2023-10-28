FROM alpine:3.18.4

COPY broker-app /app/broker
CMD ["/app/broker"]