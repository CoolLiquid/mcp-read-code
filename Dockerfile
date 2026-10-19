# Binary is built on the host (see scripts/build.sh) because golang base images
# may be unavailable on restricted networks.
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY gitlab-code-reader-mcp /usr/local/bin/gitlab-code-reader-mcp
ENV MCP_TRANSPORT=http
ENV LISTEN_ADDR=:8088
ENV GITLAB_URL=http://10.217.102.20
EXPOSE 8088
ENTRYPOINT ["gitlab-code-reader-mcp"]
