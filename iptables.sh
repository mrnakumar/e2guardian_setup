sudo iptables -I OUTPUT 1 -m owner --uid-owner root -j ACCEPT
sudo iptables -I OUTPUT 2 -p tcp -m multiport --dports 80,443 -m owner --uid-owner e2guardian -j ACCEPT
sudo iptables -I OUTPUT 3 -p tcp -m multiport --dports 80,443 -j DROP

# Restart e2guardian
sudo e2guardian -q
sleep 10
sudo e2guardian
