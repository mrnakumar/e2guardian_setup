#!/bin/bash
# Can be used to remove system wide proxies on LUbuntu

if [ $(id -u) -ne 0 ]; then
  echo "This script must be run as root";
  exit 1;
fi

gsettings set org.gnome.system.proxy mode 'none' ;

grep PATH /etc/environment > lol.t;
cat lol.t > /etc/environment;

printf "" > /etc/apt/apt.conf.d/95proxies;

rm -rf lol.t;

