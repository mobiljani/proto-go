FROM golang:1.25-alpine3.22 AS build

ARG SERVER_DIR=cmd/smoke-test/main.go

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ${SERVER_DIR}

# Runtime stage
FROM golang:1.25-alpine3.22 AS runtime

ARG PORT=8080

EXPOSE ${PORT}

RUN adduser -D runtimeuser
USER runtimeuser

# Set the workdir
WORKDIR /home/runtimeuser

COPY --from=build /server /server
# Run
CMD ["/server"]