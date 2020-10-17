FROM mysql:5.7

RUN apt-get update && \
  apt-get install -y curl

ARG ENTRYKIT_VERSION=0.4.0

RUN curl -sSfLO https://github.com/progrium/entrykit/releases/download/v${ENTRYKIT_VERSION}/entrykit_${ENTRYKIT_VERSION}_Linux_x86_64.tgz && \
  tar zxf entrykit_${ENTRYKIT_VERSION}_Linux_x86_64.tgz && \
  mv entrykit /bin/entrykit && \
  chmod +x /bin/entrykit && \
  entrykit --symlink && \
  rm entrykit_${ENTRYKIT_VERSION}_Linux_x86_64.tgz

COPY files/ /

ENV MYSQL_ALLOW_EMPTY_PASSWORD yes

ENTRYPOINT [ \
  "render", \
  "/etc/mysql/conf.d/local.cnf", \
  "--", \
  "docker-entrypoint.sh", "mysqld" \
  ]
