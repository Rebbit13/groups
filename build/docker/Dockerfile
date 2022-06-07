FROM golang:1.17

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN ls

RUN go build -o groups ./cmd/main.go
RUN go get -v github.com/rubenv/sql-migrate/...
CMD sql-migrate up -env=""postgres"" && ./groups ${FLAGS}
