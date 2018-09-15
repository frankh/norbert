FROM golang:1.11 as backend-builder
WORKDIR /go/src/github.com/frankh/norbert
RUN go get github.com/golang/dep/cmd/dep

ADD Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only
RUN go build -v ./vendor/...

ADD cmd/ ./cmd/
RUN go install -v ./cmd/...

FROM node:10 as frontend-builder
WORKDIR /frontend
RUN npm install -g yarn && chmod +x /usr/local/bin/yarn

ADD /frontend/package.json /frontend/yarn.lock ./
RUN yarn install --pure-lockfile --production=false

ADD /frontend/public ./public
ADD /frontend/src ./src

RUN yarn build

FROM golang:1.11-alpine

WORKDIR /app
ENV GO111MODULE=off
RUN apk add -U ca-certificates curl gcc libc-dev

COPY checkers/ /go/src/github.com/frankh/norbert/checkers/
COPY pkg/ /go/src/github.com/frankh/norbert/pkg/

RUN mkdir ./plugins \
  && go build -buildmode=plugin -v -o ./plugins/http.so github.com/frankh/norbert/checkers/http

COPY --from=frontend-builder /frontend/build /app/public
COPY --from=backend-builder /go/bin/norbert /usr/bin/norbert

ENTRYPOINT ["norbert"]
