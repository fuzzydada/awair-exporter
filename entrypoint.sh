#!/bin/sh

# Set user and group IDs
PUID=${PUID:-99}
PGID=${PGID:-100}

# Check if group exists
if [ -z "$(getent group ${PGID})" ]; then
  # Group doesn't exist, create it
  addgroup -g ${PGID} -S awair
  GROUP_NAME=awair
else
  # Group exists, get its name
  GROUP_NAME=$(getent group ${PGID} | cut -d: -f1)
fi

# Create user
adduser -u ${PUID} -G ${GROUP_NAME} -S -s /bin/sh awair

# Set permissions
chown -R awair:${GROUP_NAME} /config
chown awair:${GROUP_NAME} /usr/local/bin/awair-exporter

# Drop privileges and execute the main process
exec su-exec awair:${GROUP_NAME} "$@"
