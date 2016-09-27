FROM debian:jessie
MAINTAINER Seán C McCord "scm@cycoresys.com"

# Alpine Linux
#FROM gliderlabs/alpine
#RUN apk-install kamailio kamailio-extras kamailio-redis kamailio-uuid kamailio-db kamailio-json kamailio-utils

# Debian Linux
RUN apt-get update && apt-get install -y curl && apt-get clean

# Install kamailio
RUN curl -qO http://deb.kamailio.org/kamailiodebkey.gpg && \
    apt-key add kamailiodebkey.gpg && \
    echo "deb http://deb.kamailio.org/kamailio jessie main\ndeb-src http://deb.kamailio.org/kamailio jessie main" > /etc/apt/sources.list.d/kamailio.list
RUN apt-get update && apt-get install -y kamailio kamailio-json-modules kamailio-postgres-modules kamailio-redis-modules kamailio-utils-modules && apt-get clean

# Install confd (we need at least 0.7.0 for the "exists" function
RUN curl -o /usr/sbin/confd -O -J -L "https://github.com/kelseyhightower/confd/releases/download/v0.11.0/confd-0.11.0-linux-amd64"
RUN chmod +x /usr/sbin/confd

# Add logshipper
#ADD logshipper.bin /usr/bin/logshipper
#RUN chmod +x /usr/bin/logshipper

# Install config and templates
ADD confd-env /etc/confd-env
ADD config /etc/kamailio
ADD dispatcher.list /data/kamailio/dispatcher.list

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD []
