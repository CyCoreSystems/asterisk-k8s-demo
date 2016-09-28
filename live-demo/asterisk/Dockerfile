# Asterisk for Docker

FROM debian:jessie
MAINTAINER Seán C McCord "scm@cycoresys.com"

ENV ASTERISK_VER 14.0.0

# Install Asterisk
RUN apt-get update && \
   apt-get install -y autoconf build-essential libjansson-dev libxml2-dev libncurses5-dev libspeex-dev libspeexdsp-dev libgsm1-dev libsrtp0-dev uuid-dev sqlite3 libsqlite3-dev libspandsp-dev pkg-config python-dev libssl-dev openssl curl && \
   apt-get clean && \
   rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /tmp
RUN curl -o /tmp/asterisk.tar.gz http://downloads.asterisk.org/pub/telephony/asterisk/releases/asterisk-${ASTERISK_VER}.tar.gz && \
   tar xf /tmp/asterisk.tar.gz && \
   cd /tmp/asterisk-${ASTERISK_VER} && \
   ./configure --with-pjproject-bundled --with-spandsp && \
   make menuselect.makeopts && \
   menuselect/menuselect --enable EXTRA-SOUNDS-EN-GSM --disable BUILD_NATIVE --disable chan_sip menuselect.makeopts && \
   make && \
   make install && \
   rm -Rf /tmp/*

# Add sound files
#ADD sounds /var/lib/asterisk/sounds/en/na

# Add entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Install config and templates
COPY config /etc/default-asterisk

WORKDIR /
ENTRYPOINT ["/entrypoint.sh"]
CMD []
