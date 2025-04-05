# Start from the official Go image
FROM golang:1.21

# Install system dependencies
RUN apt-get update && apt-get install -y \
    libvips-dev \
    libgl1-mesa-glx \
    libgl1-mesa-dev \
    libgles2-mesa \
    libwayland-dev \
    libxkbcommon-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libglfw3 \
    libglfw3-dev \
    imagemagick \
    xvfb \
    x11-utils \
    pkg-config \
    gcc \
    g++ \
    && apt-get clean

    # Install Go
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz && \
tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz && \
ln -s /usr/local/go/bin/go /usr/bin/go

# Set Go environment
ENV PATH="/usr/local/go/bin:$PATH"

# Set working directory
WORKDIR /app

# Copy all project files
COPY . .

# Tidy up Go dependencies
RUN go mod tidy


# Build your server.go file
RUN go build server.go

RUN chmod +x src/main/main

# Expose the server port
EXPOSE 8080

CMD xvfb-run --auto-servernum --server-args="-screen 0 1024x768x24" ./server


# Start the server
CMD ["./server"]
