# FROM berkgokden/go_builder_base:v0.0.1 AS builder
FROM berkgokden/go_builder_base@sha256:d3098e69bd256a6d93fc9fa81c3256796b3512ae4d724451d3e5d1644f009063 AS builder

FROM gcr.io/distroless/static@sha256:c6d5981545ce1406d33e61434c61e9452dad93ecd8397c41e89036ef977a88f4

COPY --from=builder /app/app /app/app

ENTRYPOINT ["/app/app"]
EXPOSE 9090 9091
