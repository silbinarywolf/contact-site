FROM scratch

WORKDIR /app

# Copy files
COPY ./config.json ./config.json
COPY ./.templates ./.templates
COPY ./static ./static

# Copy server
COPY ./server ./server