FROM golang:1.9.2-alpine3.7

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache git # install git for fetching dependencies
RUN go-wrapper download    # "go get -d -v ./..."
RUN go-wrapper install     # "go install -v ./..."
RUN apk del git            # remove git from the image

CMD ["go-wrapper", "run"]  # ["app"]
