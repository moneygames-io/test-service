FROM golang:1.8

WORKDIR /go/src/app
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY Gopkg.toml .
COPY Gopkg.lock .
RUN dep ensure --vendor-only

COPY . .
RUN dep ensure
RUN go install -v ./...

CMD ["app"]
