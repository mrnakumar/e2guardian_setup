iptables -P INPUT ACCEPT
iptables -P OUTPUT ACCEPT
iptables -P FORWARD ACCEPT

iptables -F INPUT
iptables -F OUTPUT
iptables -F FORWARD

# Add rules for e2guardian
sudo iptables -I OUTPUT 1 -m owner --uid-owner root -j ACCEPT
sudo iptables -I OUTPUT 2 -p tcp -m multiport --dports 80,443 -m owner --uid-owner e2guardian -j ACCEPT
sudo iptables -I OUTPUT 3 -p tcp -m multiport --dports 80,443 -j DROP

