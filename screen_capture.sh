#!/bin/bash

function send_screenshots {
  directoryName="$1"
  ecKey="$2"
  fromAddr="$3"
  toAddr="$4"

  # compress and then send
  dirPath="${PWD}/${directoryName}"
  compressedPath="${dirPath}.tgz"
  if [ ! -f "${compressedPath}" ]; then
      tar -czf "${compressedPath}" "${dirPath}" || { echo "Failed to compress ${dirPath}"; return ; }
  fi
  echo "Attempting to encrypt ${compressedPath}"
  python "./encrypt_decrypt.py" "${ecKey}" "encrypt" "${compressedPath}"
  encryptedFilePath="${compressedPath}.ec"
  if [ ! -f "$encryptedFilePath" ]; then
    echo "Unknown error. Skipping sync of ${compressedPath}"
    return
  fi
  echo "Sending screenshots file ${encryptedFilePath} ..."
  python "./mailer.py" "${fromAddr}" "${toAddr}" "SC Records ${dirName}" "${encryptedFilePath}"
  rm -f "${encryptedFilePath}" || { echo "Failed to delete ${encryptedFilePath}, exiting."; exit; }
  rm -f "${compressedPath}" || { echo "Failed to delete ${compressedPath}, exiting."; exit; }
  rm -rf "${dirPath}" || { echo "Failed to delete ${dirPath}, exiting."; exit; }
}

function take_screenshot {
  screenshotsDirectoryPath="$1"
  if [ "${screenshotsDirectoryPath}" == "" ]; then
    echo "Require directory path to take screenshots"
    return
  fi
  currentWorkingDir="${PWD}"
  cd "${screenshotsDirectoryPath}" ||  { echo "Unknown error in cd back to screenshots directory. Returning early."; return; }
  # Take screenshot using scrot utility
  DISPLAY=:0 scrot
  cd "${currentWorkingDir}" || { echo "Unknown error in cd back to working directory";  }
}

function screenshots {
  user="$1"
  workDir="${2}/screenshots"
  ecKey="$3"
  fromAddr="$4"
  toAddr="$5"

  today=$(date +%d%m%y)
  screenShotDirPrefix="sc_"
  todaysShots="${workDir}/${screenShotDirPrefix}${today}"

  cd "${workDir}" || { echo "Failed to cd into ${workDir}. Exiting"; exit 2; }
  echo "Here at ${workDir}"
  source "./venv/bin/activate" || { echo "Could not activate python venv.";  exit; }

  # Send any existing screenshots from previous days
  for dirName in *; do
    if [ ! -d "$dirName" ]; then
      continue
    fi
    if [[ ${dirName} == ${screenShotDirPrefix}* ]]; then
      if [ "${todaysShots}" != "${workDir}/${dirName}" ]; then
        send_screenshots "${dirName}" "${ecKey}" "${fromAddr}" "${toAddr}"
      fi
    fi
  done
  if [ ! -d "${todaysShots}" ]; then
    mkdir "${todaysShots}" || { echo "Error in creating directory ${todaysShots}. Exiting"; exit 1; }
  fi
  take_screenshot "${todaysShots}"
}

function setup_cron_job(){
  user="$1"
  operationsDir="$2"
  ecKey="$3"
  fromAddr="$4"
  toAddr="$5"
  workDir="${operationsDir}/screenshots"
  sudo mkdir -p "${workDir}" || { echo "Failed to create workDir. Exiting."; exit; }
  sudo cp "./mailer.py" "./encrypt_decrypt.py" "./screen_capture.sh" "./requirements-for-python-code.txt" "./token.json" "${workDir}/" || { echo "Failed to copy data/program files. Exiting."; exit; }
  currentDir="$PWD"
  cd "${workDir}" || { echo "Failed to cd into ${workDir}. Exiting"; exit ;}
  python3 -m venv venv || { "Failed to create python virtual environment. Exiting."; exit; }
  source "./venv/bin/activate"
  pip install -r "./requirements-for-python-code.txt" || { echo "Pip install failed for venv. Exiting."; exit; }
  cd "${currentDir}"  || { echo "Failed to cd into ${currentDir}. Exiting"; exit ;}
  sudo chown "${user}":"${user}" -R "${workDir}"
  sudo chmod 700 -R "${workDir}"

  #write out current crontab
  sudo crontab -u "${user}" -l > mycron
  # Run every 3 minutes
  echo "* * * * * ${workDir}/screen_capture.sh 'send_captured' '${operationsDir}' '${user}' '${ecKey}' '${fromAddr}' '${toAddr}' > /var/log/e2guardian/screenshots.log 2>&1" >> mycron
  #install new cron file
  sudo crontab -u "${user}" mycron
  rm mycron
  echo "Successfully setup cron job for chrome history sync"
}

# Main logic starts here.
if [ "$#" -lt 5 ]; then
  echo "Usages: Operation<setup_capture|send_captured> OperationsDir EncryptionKey FromEmailAddress ToEmailAddress"
  exit 1
fi
user="e2guardian"
operation="$2"
operationsDir="$3"
ecKey="$4"
fromAddr="$5"
toAddr="$6"
workDir="${operationsDir}/screenshots"
if [ ! -d "${workDir}" ]; then
  mkdir "${workDir}" || { echo "Failed to create ${workDir}. Exiting"; exit 1; }
fi
# Validate command line arguments
if id "${user}" &>/dev/null; then
  if [ "$operation" == "setup_capture" ]; then
    echo "Will setup screenshots job for user ${user}"
  fi
else
  echo "User ${user} not found. Exiting"
  exit 3
fi
if [ ! -d "$workDir" ]; then
  echo "Work directory '${workDir}' does not exist. Exiting"
  exit 4
fi
if [ "${ecKey}" == "" ]; then
  echo "Required encryptionKey to process"
  exit 5
fi
if [ "${fromAddr}" == "" ]; then
  echo "Required sender' email address to process"
  exit 6
fi
if [ "${toAddr}" == "" ]; then
  echo "Required receiver's email address to process"
  exit 7
fi

# Call appropriate function based on requested operation
if [ "$operation" == "send_captured" ]; then
  screenshots "$user" "$operationsDir" "$ecKey" "$fromAddr" "$toAddr"
fi
if [ "$1" == "setup_capture" ]; then
  setup_cron_job "${user}" "${operationsDir}" "$ecKey" "$fromAddr" "$toAddr"
fi

