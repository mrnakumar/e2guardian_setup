#!/bin/sh
###
# This file can be used to install e2guardian.
###

sudo chmod 755 -R "/etc/e2guardian"
sudo chown "e2guardian:e2guardian" -R "/etc/e2guardian"
#sudo chown "${user}:${user}" -R "${operationsDir}"
sudo chown "e2guardian:e2guardian" -R "/var/log/e2guardian"

operationsDir="/etc/e2guardian/operations"

# Create user
if [ "$#" -lt 4 ]; then
    echo "Usages: EncryptionKey captureAndUsageUser FromEmailAddress ToEmailAddress"
    exit 1
fi
encryptionKey="$1"
user="$2"
fromEmailAddress="$3"
toEmailAddress="$4"

echo "Adding user ${user}..."
sudo useradd -m "${user}"

echo "Installing e2guardian..."
if [ ! -d "${operationsDir}" ]; then
  sudo mkdir -p "${operationsDir}" || { echo "Could not create directory ${operationsDir}. Exiting"; }
fi
sudo chmod 777 -R "${operationsDir}" || { echo "Failed to change permissions for ${operationsDir}. Exiting"; exit 1; }
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

if [ ! -d "${operationsDir}" ]; then
  sudo mkdir "${operationsDir}" || { echo "Failed to create ${operationsDir}. Exiting."; exit  1; }
fi

# Setup log rotation
#TODO: make logrotate size based so that cron can run more frequently and also fix path problems
echo "Setting up log rotation"
logRotationDir="${operationsDir}/e2guardian_log_rotate"
sudo mkdir -p "${logRotationDir}"
sudo cp "${e2GuardianDir}/data/scripts/logrotation" "${logRotationDir}/"
sudo chown e2guardian:e2guardian -R "${operationsDir}"

#write out current crontab
sudo crontab -u e2guardian -l > mycron
#echo new cron into cron file, runs on monday
echo "0 15 * * 1 ${logRotationDir}/logrotation " >> mycron
#install new cron file
sudo crontab -u e2guardian mycron
rm mycron
echo "Setting up log rotation finished successfully"




## Setup necessary services to start on boot

# Make e2guardian run on boot
sudo cp "${e2GuardianDir}/data/scripts/e2guardian.service" /etc/systemd/system/
sudo systemctl enable e2guardian

# Setup chrome sync cron job
./cron_chrome_history_sync.sh "setup_cron" "${operationsDir}" "${user}" "${encryptionKey}" "${fromEmailAddress}" "${toEmailAddress}"
cd "${e2GuardianParentDir}" || { echo "After setting up cron chrome history sync, failed to cd back into ${e2GuardianParentDir}. Exiting."; exit 1; }

# Setup screenshots cron job
./screen_capture.sh "setup_capture" "${user}" "${operationsDir}" "${encryptionKey}" "${fromEmailAddress}" "${toEmailAddress}"
cd "${e2GuardianParentDir}" || { echo "After setting up cron job for screenshots, failed to cd back into ${e2GuardianParentDir}. Exiting."; exit 1; }

# Configure iptables to make sure access is only via e2guardian
iptablesScriptPath="${operationsDir}/scripts"
if [ ! -d "${iptablesScriptPath}" ]; then
  sudo mkdir "${iptablesScriptPath}" || { echo "Could not create ${iptablesScriptPath}. Exiting"; exit 1; }
fi
echo "Copyting iptables.sh"
sudo cp "iptables.sh" "${iptablesScriptPath}/"
echo "Copied iptables.sh"
sudo cp "iptables.service" /etc/systemd/system/
sudo systemctl enable iptables.service

echo "Setting up e2guardian finished successfully"
