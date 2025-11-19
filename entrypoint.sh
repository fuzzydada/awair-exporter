#!/bin/sh

# Set user and group IDs
PUID=${PUID:-99}
PGID=${PGID:-100}

# Create group and user
groupadd -g ${PGID} -o awair
useradd --shell /bin/sh -u ${PUID} -g ${PGID} -o -c "" -m awair

# Set permissions
chown -R awair:awair /config
chown awair:awair /root/awair-exporter

# Drop privileges and execute the main process
exec su-exec awair:awair "$@"
