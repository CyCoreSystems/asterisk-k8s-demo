#!/bin/sh
set -e

# Default values
#: ${PUBLIC_IPV4:="$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)"}
#: ${PRIVATE_IPV4:="$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)"}
: ${PUBLIC_IPV4:="$(curl -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip)"}
: ${PRIVATE_IPV4:="$(curl -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip)"}
: ${RTPPROXY_ARGS:="-f -A ${PUBLIC_IPV4} -F -l ${PRIVATE_IPV4} -m 20000 -M 30000 -s udp:127.0.0.1:7722 -d INFO"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

# Run rtpproxy
exec /usr/bin/rtpproxy ${RTPPROXY_ARGS}
