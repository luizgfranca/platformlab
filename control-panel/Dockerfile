#### BACKEND
FROM golang:1.23 AS backend-build
WORKDIR /app
COPY control-panel .
RUN go build

#### FRONTEND
FROM node:22 AS frontend-build

WORKDIR /app
RUN mkdir -p jscore
COPY packages/js-core jscore

WORKDIR /app/jscore
RUN npm install
RUN npm link

WORKDIR /app
RUN mkdir -p web
WORKDIR /app/web
COPY control-panel/web/ .
RUN yarn install
RUN npm link @synthreon/core
RUN yarn build

#### FINAL CONTAINER
FROM debian:bookworm AS final
WORKDIR /app

# default sane values for some environment variables
ENV RETRY_TIMEOUT_SECONDS=60
ENV STATIC_FILES_DIR=web/dist
ENV DATABASE=/data/db.sqlite

RUN mkdir -p web/dist
COPY --from=backend-build /app/synthreon .
COPY --from=frontend-build /app/web/dist web/dist
RUN mkdir /data

CMD ["/bin/sh", "-c", "/app/controlpanel"]
