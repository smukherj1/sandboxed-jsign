#!/usr/bin/bash

set -eu

BRIDGE="docker-br0"
KMS_IP="10.0.0.2"
KMS_LOCAL_PORT="9000"
TS_IP="10.0.0.3"
TS_LOCAL_PORT="9001"

docker network rm ${BRIDGE} > /dev/null 2>&1 || echo "Ignoring error deleting non-existent docker bridge"
docker network create -d bridge ${BRIDGE} -o com.docker.network.bridge.name=${BRIDGE}
sudo iptables-save > ~/iptables.txt


sudo iptables -t nat -A PREROUTING -p tcp -i ${BRIDGE} -d ${KMS_IP} --dport 443 -j DNAT --to-destination 127.0.0.1:${KMS_LOCAL_PORT}
sudo iptables -t nat -A POSTROUTING -p tcp -o ${BRIDGE} -s 127.0.0.1 --sport ${KMS_LOCAL_PORT} -j SNAT --to-source ${KMS_IP}:443
sudo iptables -t nat -A PREROUTING -p tcp -i ${BRIDGE} -d ${TS_IP} --dport 443 -j DNAT --to-destination 127.0.0.1:${TS_LOCAL_PORT}
sudo iptables -t nat -A POSTROUTING -p tcp -o ${BRIDGE} -s 127.0.0.1 --sport ${TS_LOCAL_PORT} -j SNAT --to-source ${TS_IP}:443

sudo sysctl -w net.ipv4.conf.all.route_localnet=1
sudo sysctl net.ipv4.ip_forward=1 
# 3. Allow forwarded traffic
# iptables -A FORWARD -p tcp -s $CONTAINER_IP -d 0.0.0.0 --dport 9000 -j ACCEPT

# 4. Drop other TCP traffic from the container
# iptables -A FORWARD -i br-4acca6aa2f90 -p tcp -j DROP

# 5. Drop other UDP traffic from the container
# iptables -A FORWARD -s $CONTAINER_IP -p udp -j DROP