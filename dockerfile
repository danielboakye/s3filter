FROM golang:1.19-alpine

WORKDIR /app

COPY . /app

# Build the Go binary
RUN CGO_ENABLED=0 go build -o s3filter ./cmd

# Now you can run the executable and pass arguments at the run time.
ENTRYPOINT ["./s3filter"]