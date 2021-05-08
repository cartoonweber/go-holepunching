## We'll choose the incredibly lightweight
## Go alpine image to work withn
FROM golang:1.16.4-alpine3.13 AS builder
## We create an /app directory within our
## image that will hold our application source
## files
# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./
# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

## the lightweight scratch image we'll
## run our application within
FROM alpine:3.13.5 AS production
## We have to copy the output from our
## builder stage to our production stage
COPY --from=builder /app .
## we can then kick off our newly compiled
## binary exectuable!!
CMD ["/server"]