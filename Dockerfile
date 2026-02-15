FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Binary
COPY realtime-poll /app/realtime-poll

# Config files
COPY config/app.json /app/config/app.json
COPY config/docs/ /app/config/docs/

# Env
COPY .prod.env /app/.prod.env

EXPOSE 8080

CMD ["./realtime-poll"]
