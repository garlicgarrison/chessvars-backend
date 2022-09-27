FROM golang:1.18.3-buster AS builder

RUN apt-get update \
    && apt-get install -y  --no-install-recommends \
        git \
        make

ENV APP_DIR $GOPATH/src/github.com/garlicgarrison/chessvars-backend
WORKDIR ${APP_DIR}



COPY . ${APP_DIR}

RUN make build

FROM alpine:latest

ENV PORT=8080
EXPOSE $PORT

COPY --from=builder /go/src/github.com/garlicgarrison/chessvars-backend/bin/backend /
ENTRYPOINT [ "/backend" ]