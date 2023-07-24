FROM golang:1.16.9

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

COPY cmd/ cmd/
COPY pkg/ pkg/
#COPY vendor/ vendor/

RUN CGO_ENABLED=0 GOOS=linux go install -v /app/cmd/liqbot && \
    rm -rf *

CMD ["liqbot"]
