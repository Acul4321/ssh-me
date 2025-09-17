FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /tmp/terminal-app

FROM scratch

WORKDIR /app

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /tmp/terminal-app /app/terminal-app
# If you have migrations, uncomment the next line
# COPY --from=0 /app/migrations/ /app/migrations/

CMD ["/app/terminal-app"]