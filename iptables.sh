sudo iptables -P INPUT ACCEPT
sudo iptables -P OUTPUT ACCEPT
sudo iptables -P FORWARD ACCEPT

sudo iptables -F INPUT
sudo iptables -F OUTPUT
sudo iptables -F FORWARD

# Add rules for e2guardian
sudo iptables -I OUTPUT 1 -m owner --uid-owner root -j ACCEPT
sudo iptables -I OUTPUT 2 -p tcp -m multiport --dports 80,443 -m owner --uid-owner e2guardian -j ACCEPT
sudo iptables -I OUTPUT 3 -p tcp -m multiport --dports 80,443 -j DROP

