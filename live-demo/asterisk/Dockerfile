# Asterisk for Kubernetes
#
# It is expected that the configuration should be generated separately, as from https://github.com/CyCoreSystems/asterisk-config.
#

FROM debian:stretch
MAINTAINER Seán C McCord "ulexus@gmail.com"

ENV ASTERISK_VER 15.4.1

# Install Asterisk
RUN apt-get update && \
   apt-get install -y autoconf build-essential libjansson-dev libxml2-dev libncurses5-dev libspeex-dev libcurl4-openssl-dev libspeexdsp-dev libgsm1-dev libsrtp0-dev uuid-dev sqlite3 libsqlite3-dev libspandsp-dev pkg-config python-dev libssl-dev openssl libopus-dev liburiparser-dev xmlstarlet curl wget && \
   apt-get clean && \
   rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /tmp
RUN curl -o /tmp/asterisk.tar.gz http://downloads.asterisk.org/pub/telephony/asterisk/releases/asterisk-${ASTERISK_VER}.tar.gz && \
   tar xf /tmp/asterisk.tar.gz && \
   cd /tmp/asterisk-${ASTERISK_VER} && \
   ./configure --with-pjproject-bundled --with-spandsp --with-opus && \
   make menuselect.makeopts && \
   menuselect/menuselect --disable CORE-SOUNDS-EN-GSM --enable CORE-SOUNDS-EN-ULAW --enable EXTRA-SOUNDS-EN-ULAW --enable codec_opus --disable BUILD_NATIVE --disable chan_sip menuselect.makeopts && \
   make && \
   make install && \
   rm -Rf /tmp/*

# Add entrypoint script
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /
ENTRYPOINT ["/entrypoint.sh"]
CMD []
