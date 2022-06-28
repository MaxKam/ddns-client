#!/bin/bash

echo "Installing DDNS-Client"
echo "Please paste the full path to the GCP credentials file that DDNS-Client should use: \n"
read file_path
echo " - Updating systemd unit file for DDNS-Client with GCP credentials path"
sed -i "s|<Path to GCP authentication json file>|$file_path|" ddns-client.service

echo " - Creating ddns-client user account"
useradd -r -U -s /usr/bin/nologin ddns-client

echo " - Creating /opt/ddns-client directory"
mkdir -p /opt/ddns-client 

echo " - Copying ddns-client binary to /opt/ddns-client"
cp ./ddns-client-Linux-x64 /opt/ddns-client/

echo " - Copying ddns-client config file to /etc/"
cp ./ddns-client-config-sample.yaml /etc/ddns-client-config.yaml

echo " - Setting ddns-client group as owner of /opt/ddns-client/ directory"
chown -R :ddns-client /opt/ddns-client/

echo " - Setting ddns-client group as owner of config file in /etc"
chown -R :ddns-client /etc/ddns-client-config.yaml

echo " - Setting full permissions for ddns-client group and removing all users access"
chmod 770 /opt/ddns-client/

echo " - Copying systemd service and timer files to /etc/systemd/system/"
cp ./ddns-client.service /etc/systemd/system/
cp ./ddns-client.timer /etc/systemd/system/

echo " - enabling systemd service"
systemctl enable ddns-client.service

echo " - enabling systemd timer to run at startup and every 8 hours"
systemctl enable ddns-client.timer

echo "INSTALL COMPLETE. Please edit config file located at /etc/ddns-client-config.yaml."
echo "Then either reboot machine or run this command to start the service: sudo systemctl start ddns-client.service"


