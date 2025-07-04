FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES git build-base

# Set working directory for the build
WORKDIR /go/src/github.com/evmos/ethermint

# Install dependencies
RUN apk add --update $PACKAGES
RUN apk add linux-headers

# Add source files
COPY . .

# Make the binary
RUN make build

# Final image
FROM alpine:3.17.3

# Install ca-certificates
RUN apk add --update ca-certificates jq
WORKDIR /

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/evmos/ethermint/build/monchaind /usr/bin/monchaind

# Run monchaind by default
CMD ["monchaind"]
