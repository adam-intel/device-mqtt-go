
#
# Copyright (c) 2018, 2019 Intel
#
# SPDX-License-Identifier: Apache-2.0
#
FROM golang:1.11-alpine AS builder

ENV GO111MODULE=on
WORKDIR /device-mqtt-go

LABEL license='SPDX-License-Identifier: Apache-2.0' \
    copyright='Copyright (c) 2018, 2019: Intel'

RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

# add git for go modules
RUN apk update && apk add make git

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o device-mqtt-go cmd/main.go

# Next image - Copy built Go binary into new workspace
FROM alpine

LABEL license='SPDX-License-Identifier: Apache-2.0' \
    copyright='Copyright (c) 2018, 2019: Intel'

ENV APP_PORT=49982
#expose command data port
EXPOSE $APP_PORT

WORKDIR /
COPY --from=builder /device-mqtt-go/device-mqtt-go /usr/local/bin/device-mqtt-go
COPY --from=builder /device-mqtt-go/cmd/res/ /res/

CMD [ "/usr/local/bin/device-mqtt-go","--profile=docker","--confdir=/res"]