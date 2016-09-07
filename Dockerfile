FROM alpine:3.3

COPY .dist/tugbot-kubernetes /usr/bin/tugbot-kubernetes

LABEL tugbot=kubernetes

ENTRYPOINT ["/usr/bin/tugbot-kubernetes"]