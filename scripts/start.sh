#!/bin/bash -e

echo "clamsig-puller v0.0.4"

# This is useful so we can debug containers running inside of OpenShift that are
# failing to start properly.

if [ "$OO_PAUSE_ON_START" = "true" ] ; then
  echo
  echo "This container's startup has been paused indefinitely because OO_PAUSE_ON_START has been set."
  echo
  while true; do
    sleep 10
  done
fi

if [ "$INIT_CONTAINER" = "true" ] ; then
  echo
  echo "The INIT_CONTAINER variable has been set. This container will attempt to populate the pod's shared volume with clam DBs and config files before exiting."
  echo
  /bin/clamsig-puller
  echo "Finished pulling signatures and config files - exiting."
  exit 0
fi

echo This container hosts the following applications:
echo
echo '/usr/bin/clamsig-puller'
echo
echo 'Every 12 hours, check if there are any databases in our mirror that are newer than what is already on disk.'
echo '----------------'
/usr/local/bin/loop.sh 43200 /bin/clamsig-puller
