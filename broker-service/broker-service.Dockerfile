FROM alpine:3.18.4

COPY app /app/broker
CMD ["/app/broker"]