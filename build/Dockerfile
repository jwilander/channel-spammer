ARG DOCKER_BUILD_IMAGE=golang:1.14
ARG DOCKER_BASE_IMAGE=alpine:3.12

FROM ${DOCKER_BUILD_IMAGE} AS build
WORKDIR /channel-spammer/
COPY . /channel-spammer/
RUN make build

# Final Image
FROM ${DOCKER_BASE_IMAGE}
COPY words /usr/share/dict/words

COPY --from=build /channel-spammer/build/_output/bin/main /channel-spammer/main
WORKDIR /channel-spammer/

CMD ./main
