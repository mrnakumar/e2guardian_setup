#!/bin/sh
sudo mkdir -p /usr/local/share/e2guardian_log_rotate
sudo cp ./data/scripts/logrotation /usr/local/share/e2guardian_log_rotate/
sudo chown root:root -R /usr/local/share/e2guardian_log_rotate
#write out current crontab
sudo crontab -u root -l > mycron
#echo new cron into cron file, runs on monday
echo "0 15 * * 1 /usr/local/share/e2guardian_log_rotate/logrotation " >> mycron
#install new cron file
sudo crontab -u root mycron
rm mycron
