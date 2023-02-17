FROM dzo.sw.sbc.space/amazme_dev/ci90000188_amazme_dev/postgres:13.6

RUN apt-get update && apt-get install -y curl bash

ENV PGPASSWORD=password
ENV PGUSER=admin
ENV PGHOST=127.0.0.1
ENV PGDATABASE=postgres

WORKDIR /app

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

COPY scripts/*.sh /app/
COPY db/migration /app/migration

RUN chmod a+x *.sh

USER 1001

CMD ["sh", "-c", "/app/start_with_migrations.sh"]