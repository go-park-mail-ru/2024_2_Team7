FROM debian:bookworm-slim
WORKDIR /app
COPY ./bin/auth_service /auth_service_run
EXPOSE 50051
CMD ["/auth_service_run"]
