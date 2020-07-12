FROM scratch

WORKDIR /app

# Copy assets files
COPY ./config.json ./config.json
COPY ./.templates ./.templates
COPY ./static ./static

# Copy server binary file
COPY ./server ./server
