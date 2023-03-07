FROM golang:1.20 as base

RUN apt-get update && apt-get install -y inetutils-ping

WORKDIR /app

COPY . . 

RUN go mod tidy 

FROM base as api

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD [ "sh", "-c", "curl --user-agent 'Internal Container HealthCheck' -H 'Content-Type: text/plain' -f http://localhost:3000/health || exit 1;" ]
