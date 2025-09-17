FROM denoland/deno:debian-2.4.3

EXPOSE 8000

WORKDIR /app

# These steps will be re-run upon each file change in your working directory:
COPY . .

# Change the ownership of the files to the deno user
RUN chown -R deno:deno .

USER deno

# Compile the main app so that it doesn't need to be compiled each startup/entry.
RUN deno cache src/main.ts

CMD ["/bin/sh", "-c", "deno task build && deno task serve"]