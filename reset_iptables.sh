#!/bin/bash
sudo iptables -P INPUT ACCEPT
sudo iptables -P OUTPUT ACCEPT
sudo iptables -P FORWARD ACCEPT

sudo iptables -F INPUT
sudo iptables -F OUTPUT
sudo iptables -F FORWARD

