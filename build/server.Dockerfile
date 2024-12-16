FROM debian:bookworm-slim
WORKDIR /app
COPY ../bin/server_service /server_service_run
COPY ./static /static
EXPOSE 8080
CMD ["/server_service_run"]