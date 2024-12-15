FROM debian:bookworm-slim
WORKDIR /app
COPY ./bin/image_service /image_service_run
COPY ./static /static
EXPOSE 50054
CMD ["/image_service_run"]