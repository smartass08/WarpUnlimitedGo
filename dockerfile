FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64 \
    TZ=Asia/Kolkata

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Command to run when starting the container
CMD ["/build/main"]
