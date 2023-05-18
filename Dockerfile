FROM golang:1.20-alpine3.16 AS binarybuilder
USER root

ENV DIR=/go/hrp

RUN go env -w GO111MODULE=on; \
    mkdir -p ${DIR}; \
    apk add --update --no-cache make

COPY . ${DIR}

WORKDIR ${DIR}

RUN make build

FROM alpine:3.16
USER root

ENV DIR=/go/hrp

RUN apk add --update --no-cache tzdata

COPY --from=binarybuilder ${DIR}/bin/ /usr/bin/

EXPOSE 80

CMD ["hrp"]