FROM golang:1.13.4-stretch AS build

WORKDIR /usr/src/tech-db

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
#RUN make build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOGC=off go build -a -installsuffix cgo -ldflags="-w -s" -v -o ./technopark-db-forum ./cmd/technopark-db-forum

FROM ubuntu:18.04 AS release

MAINTAINER Nozim Yunusov

#
# Установка postgresql
#
ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "include_dir='conf.d'" >> /etc/postgresql/$PGVER/main/postgresql.conf
ADD ./postgresql.conf /etc/postgresql/$PGVER/main/conf.d/basic.conf

#RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "shared_buffers = 256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "wal_buffers = 32MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "work_mem = 32MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "maintenance_work_mem = 320MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "huge_pages = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "effective_cache_size = 512MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_error_verbosity = TERSE" >> /etc/postgresql/$PGVER/main/postgresql.conf



# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Объявлем порт сервера
EXPOSE 5000

COPY ./assets/db/postgres/base.sql ./assets/db/postgres/base.sql
# Собранный ранее сервер
COPY --from=build /usr/src/tech-db/technopark-db-forum .

#
# Запускаем PostgreSQL и сервер
#
CMD service postgresql start && ./technopark-db-forum