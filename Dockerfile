FROM golang:1.11-alpine as backend-builder
ENV GO111MODULE=off
WORKDIR /go/src/github.com/frankh/norbert
RUN apk add -U git gcc libc-dev

RUN go get github.com/golang/dep/cmd/dep
RUN go get github.com/gobuffalo/packr/...

ADD Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only
RUN go build -v ./vendor/...

ADD pkg/ ./pkg/
RUN go build -v ./pkg/...

ADD cmd/ ./cmd/
RUN packr
RUN go install -v ./cmd/...

FROM node:10 as frontend-builder
WORKDIR /frontend
RUN npm install -g yarn && chmod +x /usr/local/bin/yarn

ADD /frontend/package.json /frontend/yarn.lock ./
RUN yarn install --pure-lockfile --production=false

ADD /frontend/public ./public
ADD /frontend/src ./src

RUN yarn build

# Final image
FROM golang:1.11-alpine

WORKDIR /app
ENV GO111MODULE=off
RUN apk add -U ca-certificates curl gcc libc-dev

COPY checkrunners/ /go/src/github.com/frankh/norbert/checkrunners/
COPY pkg/ /go/src/github.com/frankh/norbert/pkg/

# Warmup plugin build caches
RUN mkdir ./plugins \
  && go build -buildmode=plugin -v github.com/frankh/norbert/checkrunners/http

COPY --from=frontend-builder /frontend/build /app/public
COPY --from=backend-builder /go/bin/norbert /usr/bin/norbert

ENTRYPOINT ["norbert"]
