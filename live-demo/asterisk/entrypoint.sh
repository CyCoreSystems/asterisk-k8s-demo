#!/bin/sh
set -e

# Default values
: ${ASTERISK_ARGS:="-fpvvvvvv"}
#: ${PUBLIC_IPV4:="$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)"}
: ${PUBLIC_IPV4:="$(curl -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip)"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

# Copy the default Asterisk configuration into place
# on the config volume
mkdir -p /etc/asterisk
cp -an /etc/default-asterisk/* /etc/asterisk/

# Build the environment-dependant configurations
mkdir -p /etc/asterisk/pjsip.d
cat <<ENDHERE >/etc/asterisk/pjsip.d/transport-udp.conf
[transport-udp]
type=transport
protocol=udp
bind=0.0.0.0
;local_net=10.0.0.0/8
external_media_address=${PUBLIC_IPV4}
external_signaling_address=${PUBLIC_IPV4}
ENDHERE

# Run Asterisk
exec /usr/sbin/asterisk ${ASTERISK_ARGS}
