# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.16-alpine

# Add Maintainer Info
LABEL maintainer="Thisara Pramuditha <thisapramuditha@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o exchange-service .

# Expose port 9090 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./exchange-service"]