FROM gitlab.amazmetest.ru:5050/devops/docker/amazme/python:3.9-slim-buster

ENV KAFKA_BROKERS   "127.0.0.1:9092"
ENV KAFKA_USER      "alice"
ENV KAFKA_PASSWORD  "alice-secret"
ENV KAFKA_SSL_CERT  ""

WORKDIR /app

COPY scripts/topics_creator/requirements.txt requirements.txt

RUN pip3 install -r requirements.txt

COPY scripts/topics_creator .