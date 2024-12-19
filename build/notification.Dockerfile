FROM debian:bookworm-slim
WORKDIR /app
COPY ./bin/notification_service /notification_service_run
EXPOSE 50056
CMD ["/notification_service_run"]