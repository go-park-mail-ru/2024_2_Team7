FROM debian:bookworm-slim
WORKDIR /app
COPY ./bin/event_service /event_service_run
EXPOSE 50053
CMD ["/event_service_run"]