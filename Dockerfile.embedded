FROM node:lts as ui-builder

### Copy Vue.js code
WORKDIR /app
COPY frontend frontend

### Build Vue.js into plain HTML/CSS/JS
WORKDIR /app/frontend
RUN npm run build-embed

FROM golang:1.19.1 AS builder

### Copy Go code
WORKDIR /src/
COPY . .
COPY --from=ui-builder /app/frontend frontend

### Build executables
RUN go build -tags webui -o /app/webapi ./cmd/webapi


### Create final container
FROM debian:bullseye

RUN apt update && apt install -y docker docker-compose openssh-client
RUN apt clean

### Inform Docker about which port is used
EXPOSE 3000

### Copy the build executable from the builder image
WORKDIR /app/
COPY --from=builder /app/webapi ./

### Executable command
CMD ["/app/webapi", "--db-filename", "/config/wasaphoto.db", "--data-keys-path", "/config"]