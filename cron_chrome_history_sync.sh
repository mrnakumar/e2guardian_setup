#!/bin/bash

# Script to setup cron job to schedule chrome history sync.
# Should be run once on a computer.
# If script finishes successfully then, doing: 'sudo crontab -u user -l' , should show the expected cron job
# Required:
# Log directory /var/log/e2guardian should have permission 700.And should have file cron_chrome_sync.log.
# To check for correctness, change the cron expression to run every minutes (i.e. '* * * * *'). Revert once checked.

user="guest"
workDir="/etc/${user}/chrome_history"
sudo mkdir -p "${workDir}" || { echo "Failed to create workDir. Exiting."; exit; }
sudo cp "./mailer.py" "./encrypt_decrypt.py" "./chrome-history-sync.sh" "./requirements-for-python-code.txt" "./token.json" "${workDir}/" || { echo "Failed to copy data/program files. Exiting."; exit; }
currentDir="$PWD"
cd "${workDir}" || { echo "Failed to cd into ${workDir}. Exiting"; exit ;}
python3 -m venv venv || { "Failed to create python virtual enviornment. Exiting."; exit; }
source "./venv/bin/activate"
pip install -r "./requirements-for-python-code.txt" || { echo "Pip install failed for venv. Exiting."; exit; }
cd "${currentDir}"  || { echo "Failed to cd into ${currentDir}. Exiting"; exit ;}
sudo chown "${user}":"${user}" -R "${workDir}"
sudo chmod 700 -R "${workDir}"

#write out current crontab
sudo crontab -u "${user}" -l > mycron
# Run every 20 minutes
echo "*/20 * * * * ${workDir}/chrome-history-sync.sh '${user}' '{encryptKey}' '{fromEmail}' '{toEmail}' > /var/log/e2guardian/cron_chrome_sync.log 2>&1" >> mycron
#install new cron file
sudo crontab -u "${user}" mycron
rm mycron

