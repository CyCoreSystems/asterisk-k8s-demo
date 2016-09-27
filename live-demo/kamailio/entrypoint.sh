#!/bin/bash
# Set default settings, pull repository, build
# app, etc., _if_ we are not given a different
# command.  If so, execute that command instead.
set -e

# Default values
: ${PID_FILE:="/var/run/kamailio.pid"}
: ${KAMAILIO_ARGS:="-DD -E -f /etc/kamailio/kamailio.cfg -P ${PID_FILE}"}

# confd requires that these variables actually be exported
export PID_FILE

# Make dispatcher.list exists
mkdir -p /data/kamailio
touch /data/kamailio/dispatcher.list

# Obtain private and public IPs
export PRIVATE_IPV4=$(curl http://169.254.169.254/latest/meta-data/local-ipv4)
export PUBLIC_IPV4=$(curl http://169.254.169.254/latest/meta-data/public-ipv4)
export PUBLIC_HOSTNAME=$(curl http://169.254.169.254/latest/meta-data/public-hostname)

# Build run-time configuration
confd -onetime -backend env -confdir=/etc/confd-env -config-file=/etc/confd-env/conf.d/kamailio.toml

# Runs kamaillio, while shipping stderr/stdout to logstash
exec /usr/sbin/kamailio $KAMAILIO_ARGS $*

