#!/bin/bash


echo "removing prometheus"
sudo docker stop prometheus
sudo docker rm -v prometheus

echo "removing node_exporter"
sudo docker stop node_exporter
sudo docker rm -v node_exporter

echo "removing grafana"
sudo docker stop grafana
sudo docker rm -v grafana