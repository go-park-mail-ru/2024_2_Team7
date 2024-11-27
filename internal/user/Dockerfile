FROM debian:bookworm-slim
WORKDIR /app
COPY ./bin/user_service /user_service_run
EXPOSE 50052
CMD ["/user_service_run"]