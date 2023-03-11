FROM node:lts-alpine3.17 as base

run npm i -g pnpm

WORKDIR /app

COPY package.json package.json

FROM base as dev

