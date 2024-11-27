FROM debian:bookworm-slim
WORKDIR /app
COPY ./bin/csat_service /csat_service_run
EXPOSE 50055
CMD ["/csat_service_run"]