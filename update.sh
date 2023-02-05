#!/bin/bash
set -e
docker pull mpkondrashin/remindme
docker stop remindme
docker rm remindme
docker run --name remindme -d -p 443:443 --mount src="$(pwd)",target=/db,type=bind mpkondrashin/remindme