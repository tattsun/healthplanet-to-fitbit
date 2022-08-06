FROM golang:1.19 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make all


FROM golang:1.19

WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/bin/* /usr/local/bin/

CMD ["healthplanet-to-fitbit"]
