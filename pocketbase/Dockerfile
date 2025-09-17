FROM alpine:latest

ARG PB_VERSION=0.30.0

RUN apk add --no-cache \
    unzip \
    ca-certificates

# download and unzip PocketBase
ADD https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip /tmp/pb.zip
RUN unzip /tmp/pb.zip -d /pocketbase/

# uncomment to copy the local pb_migrations dir into the image
COPY ./pb_migrations /pocketbase/pb_migrations

# uncomment to copy the local pb_hooks dir into the image
# COPY ./pb_hooks /pocketbase/pb_hooks

EXPOSE 8080

# create admin user acording to the env file
# start PocketBase
CMD ["sh", "-c", "/pocketbase/pocketbase superuser upsert \"$POCKETBASE_ADMIN_EMAIL\" \"$POCKETBASE_ADMIN_PASSWORD\" && exec /pocketbase/pocketbase serve --http=0.0.0.0:8080"]