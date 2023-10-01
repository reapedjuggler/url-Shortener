FROM golang:latest
WORKDIR /app/url-shortener
COPY . .
RUN go mod download
RUN go build -o main .  
EXPOSE 8000
CMD ["go", "run", "main.go"]