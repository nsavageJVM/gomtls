FROM golang:alpine

# docker run -p 8443:8443 go-dock

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


RUN apk add openssl nano

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
#COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .
RUN mkdir -p /build/tls
COPY tls /build/tls

# Build the application
RUN go build -o server .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/server .
RUN mkdir -p /dist/tls
RUN cp /build/tls/* /dist/tls
RUN mkdir -p /dist/tls_2
COPY tls_2 /dist/tls_2/

# Export necessary port
EXPOSE 8443

# Command to run when starting the container
CMD ["/dist/server"]