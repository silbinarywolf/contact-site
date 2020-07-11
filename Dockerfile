FROM scratch

WORKDIR /app

COPY ./server ./server
COPY ./.assets ./.assets
COPY ./static ./static
