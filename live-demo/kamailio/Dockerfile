FROM debian:stretch
MAINTAINER Seán C McCord "scm@cycoresys.com"

RUN apt-get update && \
    apt-get install -y curl procps gnupg2 && \
    apt-get clean

# Install kamailio
RUN curl -qO http://deb.kamailio.org/kamailiodebkey.gpg && \
    apt-key add kamailiodebkey.gpg && \
    echo "deb http://deb.kamailio.org/kamailio51 stretch main\ndeb-src http://deb.kamailio.org/kamailio50 stretch main" > /etc/apt/sources.list.d/kamailio.list
RUN apt-get update &&  apt-get install -y kamailio kamailio-json-modules kamailio-utils-modules kamailio-extra-modules kamailio-nth && apt-get clean

# Download netdiscover
RUN curl -qL -o /usr/bin/netdiscover https://github.com/CyCoreSystems/netdiscover/releases/download/v1.2.3/netdiscover.linux.amd64
RUN chmod +x /usr/bin/netdiscover

# Install config and templates
ADD config /etc/kamailio
ADD dispatcher.list /data/kamailio/dispatcher.list

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD []
