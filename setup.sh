#!/bin/sh
###
# This file can be used to install e2guardian.
###

# Create user
echo "Adding user e2guardian..."
useradd -m e2guardian

echo "Installing e2guardian..."

# Install requirements for e2guardian
while IFS='\n' read -r package; do
   sudo apt-get install $package -y || continue
done < requirements

# Download e2guardian
wget https://github.com/e2guardian/e2guardian/archive/refs/heads/v5.4.zip || exit

unzip "v5.4.zip" || { echo "Failed to unzip e2guardian. Exiting"; exit 1; }

# Build e2guardian
e2GuardianParentDir="${PWD}"
e2GuardianDir="${PWD}/e2guardian-5.4"
cd "${e2GuardianDir}" || { echo "Failed to cd into ${e2GuardianDir}. Exiting"; exit 1; }
echo "Compiling e2guardian..."
sudo "./autogen.sh" || exit
sudo "./configure" '--prefix=/usr' '--enable-clamd=yes' '--with-proxyuser=e2guardian' '--with-proxygroup=e2guardian' '--sysconfdir=/etc' '--localstatedir=/var' '--enable-icap=yes' '--enable-commandline=yes' '--enable-email=yes' '--enable-ntlm=yes' '--mandir=${prefix}/share/man' '--infodir=${prefix}/share/info' '--enable-pcre=yes' '--enable-sslmitm=yes' 'CPPFLAGS=-mno-sse2 -g -O2' || exit

sudo make || { echo "Failed to make. Exiting."; exit 1; }
sudo make install || { echo "Failed to make install. Exiting."; exit 1; }

cd "${e2GuardianParentDir}" || { echo "Failed to cd back into ${e2GuardianParentDir}. Exiting."; exit 1; }

# Copy files for group1
echo "Setting up group1"
./group1_setup.sh


## Operations stuff

# Make e2guardian run on boot
sudo cp "${e2GuardianDir}/data/scripts/e2guardian.service" /etc/systemd/system/
sudo systemctl enable e2guardian


# Setup log rotation
#TODO: make logrotate size based so that cron can run more frequently and also fix path problems
echo "Setting up log rotation"
sudo mkdir -p /usr/local/share/e2guardian_log_rotate
logRotationDir="/usr/local/share/e2guardian_log_rotate"
sudo cp "${e2GuardianDir}/data/scripts/logrotation" "${logRotationDir}/"
sudo chown root:root -R "${logRotationDir}"
#write out current crontab
sudo crontab -u root -l > mycron
#echo new cron into cron file, runs on monday
echo "0 15 * * 1 ${logRotationDir}/logrotation " >> mycron
#install new cron file
sudo crontab -u root mycron
rm mycron
echo "Setting up log rotation finished successfully"


# Configure iptables to make sure access is only via e2guardian
# TODO: setup iptables on boot. Config file can be copied from the other computer.
./iptables.sh
