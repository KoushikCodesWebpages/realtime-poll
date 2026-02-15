# Use a lightweight Debian-based image
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /app

# Install CA certificates (required for HTTPS/TLS, e.g., MongoDB Atlas)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the pre-built Go binary named "JSE" into the container
COPY realtime-poll /app/realtime-poll

# Copy the public folder (for serving static frontend assets)
# COPY public /app/public

# COPY app/templates /app/app/templates

# Copy the .env file (environment configuration)
COPY .prod.env /app/.prod.env

# Expose the port your app will listen on
EXPOSE 8080

# Run the Go app
CMD ["./realtime-poll"]
