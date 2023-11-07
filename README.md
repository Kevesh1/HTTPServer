# RUN Docker

## Create container (if not already inited)

docker build --tag http-server .

## Start image based on container

docker run --publish 8080:8080 http-server

We need the publish flag to connect the image port 8080 to our local 8080 port so we can access.