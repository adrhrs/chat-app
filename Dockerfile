FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go mod tidy
RUN go mod vendor
# RUN go install -v ./...

# EXPOSE 8080
RUN go build -o main .

# CMD ["mkdir -p test"] 
# RUN mkdir -p result

COPY /static app/static

# CMD ["chat-app"]
ENTRYPOINT ["/app/main"]