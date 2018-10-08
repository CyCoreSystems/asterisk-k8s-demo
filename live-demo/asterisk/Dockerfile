# Asterisk for Kubernetes
#
# It is expected that the configuration should be generated separately, as from https://github.com/CyCoreSystems/asterisk-config.
#

FROM debian:stretch as builder
MAINTAINER Seán C McCord "ulexus@gmail.com"

ENV ASTERISK_VER 15.6.1

# Install Asterisk
RUN apt-get update && \
   apt-get install -y autoconf build-essential libjansson-dev libxml2-dev libncurses5-dev libspeex-dev libcurl4-openssl-dev libspeexdsp-dev libgsm1-dev libsrtp0-dev uuid-dev sqlite3 libsqlite3-dev libspandsp-dev pkg-config python-dev libssl-dev openssl libopus-dev liburiparser-dev xmlstarlet curl wget && \
   apt-get clean && \
   rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /tmp
RUN curl -o /tmp/asterisk.tar.gz http://downloads.asterisk.org/pub/telephony/asterisk/releases/asterisk-${ASTERISK_VER}.tar.gz && \
   tar xf /tmp/asterisk.tar.gz && \
   cd /tmp/asterisk-${ASTERISK_VER} && \
   curl -L -o apps/app_audiosocket.c https://raw.githubusercontent.com/CyCoreSystems/audiosocket/master/asterisk/app_audiosocket.c && \
   ./configure --with-pjproject-bundled --with-spandsp --with-opus && \
   make menuselect.makeopts && \
   menuselect/menuselect --disable CORE-SOUNDS-EN-GSM --enable CORE-SOUNDS-EN-ULAW --enable codec_opus --disable BUILD_NATIVE --disable chan_sip menuselect.makeopts && \
   make && \
   make install && \
   rm -Rf /tmp/*

FROM debian:stretch
COPY --from=builder /usr/sbin/asterisk /usr/sbin/
COPY --from=builder /usr/sbin/safe_asterisk /usr/sbin/
COPY --from=builder /usr/lib/libasterisk* /usr/lib/
COPY --from=builder /usr/lib/asterisk/ /usr/lib/asterisk
COPY --from=builder /var/lib/asterisk/ /var/lib/asterisk
COPY --from=builder /var/spool/asterisk/ /var/spool/asterisk

# Add required runtime libs
RUN apt-get update && \
   apt-get install -y gnupg libjansson4 xml2 libncurses5 libspeex1 libcurl4-openssl-dev libspeexdsp1 libgsm1 libsrtp0 uuid libsqlite3-0 libspandsp2 libssl1.1 libopus0 liburiparser1 xmlstarlet curl wget && \
   apt-get clean && \
   rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Add sngrep
COPY irontec.list /etc/apt/sources.list.d/irontec.list
RUN curl -L http://packages.irontec.com/public.key | apt-key add -
RUN apt-get update && \
   apt-get install -y sngrep && \
   rm -Rf /var/lib/apt/lists/ /tmp/* /var/tmp/*

# Add entrypoint script
ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /
ENTRYPOINT ["/entrypoint.sh"]
CMD []
