FROM golang:alpine3.16
WORKDIR /app

# add some necessary packages
RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

# Copy and build the app
# COPY . /app
COPY ./entrypoint.sh /entrypoint.sh

# prevent the re-installation of vendors at every change in the source code
COPY ./go.mod go.sum ./
RUN go mod download && go mod verify

# Install Compile Daemon for go. We'll use it to watch changes in go files
RUN go get github.com/githubnemo/CompileDaemon

COPY . .

# RUN go build -o main main.go

# wait-for-it requires bash, which alpine doesn't ship with by default. Use wait-for instead
ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

# ENTRYPOINT [ "go","run", "main.go" ]

ENTRYPOINT [ "sh","entrypoint.sh" ]