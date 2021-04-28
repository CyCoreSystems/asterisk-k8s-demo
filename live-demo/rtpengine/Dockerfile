FROM alpine
RUN apk add --no-cache rtpengine curl

# Download netdiscover
RUN curl -qL -o /usr/bin/netdiscover https://github.com/CyCoreSystems/netdiscover/releases/download/v1.2.3/netdiscover.linux.amd64
RUN chmod +x /usr/bin/netdiscover

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD []
