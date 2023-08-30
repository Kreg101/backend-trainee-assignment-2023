FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres executable
RUN chmod +x wait-for-postgres.sh

# build go server
RUN go mod download
RUN go build -o server ./cmd/main.go

CMD ["./server"]