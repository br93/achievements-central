FROM alpine:3.18.4

COPY app /app/acc-management
CMD ["/app/acc-management"]