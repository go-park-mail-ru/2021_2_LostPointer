#!/bin/bash

#echo "removing all containers"
# shellcheck disable=SC2046
#sudo docker rm -vf $(docker ps -a -q)

echo "running prometheus on port 9090"
# shellcheck disable=SC2154
sudo docker run -p 9090:9090 -d --name prometheus --net=host -v $PWD/prometheus:/etc/config prom/prometheus --config.file=/etc/config/prometheus.yml
echo "running node exporter on port 9100"
sudo docker run -p 9100:9100 -d --name node_exporter --net=host -v $PWD/node_exporter:/etc/config prom/node-exporter --path.rootfs=/etc/config
echo "running grafana on port 3000"
sudo docker run -d -p 3000:3000 --name grafana --net=host -v $PWD/grafana/grafana.ini:/etc/grafana/grafana.ini grafana/grafana
