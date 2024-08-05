#!/bin/bash


# is docker exist?
if [ "$(command -v docker)" == "" ];
then
	echo "FATAL: you must install docker...";
	exit 1;
fi


# is docker a new?
if [ "$(docker 2>&1 | grep compose | wc -l)" == "0" ];
then
	echo "FATAL: you must upgrade docker, docker must be with compose function...";
	exit 1;
fi


# down containers
echo "down..."
docker compose --project-name pow -f ./docker-compose.yml down --remove-orphans -t 30 &> /dev/null;
if [ "${?}" != "0" ];
then
	echo "ERROR: docker compose down";
	exit 1;
fi


# build containers
echo "build..."
docker compose --project-name pow -f ./docker-compose.yml build --no-cache;
if [ "${?}" != "0" ];
then
	echo "ERROR: docker compose build";
	exit 1;
fi


# up containers
echo "up..."
docker compose --project-name pow -f ./docker-compose.yml up --renew-anon-volumes --always-recreate-deps -d &> /dev/null;
if [ "${?}" != "0" ];
then
	echo "ERROR: docker compose up";
	exit 1;
fi


# wait when container will be down
while true;
do
	if [ $(docker inspect -f '{{.State.Running}}' pow-client) = "false" ];
	then
		break;
	fi
	sleep 1; # wait
done


# show log from container
docker logs pow-client;
if [ "${?}" != "0" ];
then
	echo "ERROR: docker logs";
	exit 1;
fi


exit 0;
