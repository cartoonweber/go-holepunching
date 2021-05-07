## We'll choose the incredibly lightweight
## Go alpine image to work withn
FROM golang:1.16.4-alpine3.13 AS builder
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
ADD . /app
## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app
## we run go build to compile the binary
## executable of our Go program
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...

## the lightweight scratch image we'll
## run our application within
FROM alpine:3.13.5 AS production
## We have to copy the output from our
## builder stage to our production stage
COPY --from=builder /app .
## we can then kick off our newly compiled
## binary exectuable!!
CMD ["./main"]