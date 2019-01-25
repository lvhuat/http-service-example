FROM ubuntu:16.04
RUN mkdir -p /app
COPY user app.toml /app/
WORKDIR /app
CMD ["./user"]

# http service
EXPOSE 8080