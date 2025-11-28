FROM ubuntu:latest
WORKDIR /app
ENV DB_HOST=localhost DB_PORT=5432
ENV DB_USER=root DB_PASSWORD=root DB_NAME=root
COPY ./main main
EXPOSE 8000
CMD [ "./main" ]