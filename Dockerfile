FROM golang:1.11 as backend-builder
WORKDIR /app

ADD go.mod go.sum ./
RUN go mod download

ADD cmd/ ./cmd/
RUN CGO_ENABLED=0 go install -v ./cmd/...

FROM node:10 as frontend-builder
WORKDIR /frontend
RUN npm install -g yarn && chmod +x /usr/local/bin/yarn

ADD /frontend/package.json /frontend/yarn.lock ./
RUN yarn install --pure-lockfile --production=false

ADD /frontend/public ./public
ADD /frontend/src ./src

RUN yarn build

FROM alpine

WORKDIR /app
RUN apk add -U ca-certificates curl
RUN mkdir -p /data

COPY --from=frontend-builder /frontend/build /app/public
COPY --from=backend-builder /go/bin/norbert /usr/bin/norbert

ENTRYPOINT ["norbert"]
