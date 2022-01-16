#!/bin/bash

# Based on requested operation (via command line argument), this script can do either:
# 1. Setup cron job to send chrome history periodically
# 2. Send chrome history
# The setting up of cron job should be done only once on a computer.
# The cron job will invoke the send chrome history function periodically thereafter.
# If cron job setup finishes successfully then, doing: 'sudo crontab -u user -l' , should show the expected cron job
#
# Required:
# Log directory /var/log/e2guardian should have permission 700.And should have file cron_chrome_sync.log.
# To check for correctness, change the cron expression to run every minutes (i.e. '* * * * *'). Revert once checked.

function setup_cron_job(){
    operationsDir="$1"
    user="$2"
    ecKey="$3"
    fromAddr="$4"
    toAddr="$5"
    if id "${user}" &>/dev/null; then
        echo "Will setup chrome sync cron job for user ${user}"
    else
        echo "User ${user} not found. Exiting"
	      exit 3
    fi

    workDir="${operationsDir}/chrome_history"
    sudo mkdir -p "${workDir}" || { echo "Failed to create workDir. Exiting."; exit; }
    sudo chmod 777 -R "${workDir}"
    sudo cp "./mailer.py" "./encrypt_decrypt.py" "./cron_chrome_history_sync.sh" "./requirements-for-python-code.txt" "./token.json" "${workDir}/" || { echo "Failed to copy data/program files. Exiting."; exit; }
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
    echo "*/20 * * * * ${workDir}/cron_chrome_history_sync.sh 'sync_chrome' '${operationsDir}' '${user}' '${ecKey}' '${fromAddr}' '${toAddr}' > /var/log/e2guardian/cron_chrome_sync.log 2>&1" >> mycron
    #install new cron file
    sudo crontab -u "${user}" mycron
    rm mycron
    echo "Successfully setup cron job for chrome history sync"
}

##########

# Sync chrome usages. Needs sqlite3 in path
function sync_chrome_history() {
    if [ "$#" -ne 5 ] ; then
        echo "Usage: $0 OperationsDir UserName EncryptionKey FromEmailAddress ToEmailAddress"
        exit 1
    fi

    # user is used to construct path to chrome History file and workDir
    operationsDir="$1"
    user="$2"
    ecKey="$3"
    fromAddr="$4"
    toAddr="$5"

    LOCKFILE="/tmp/chrome_sync_lock.txt"
    if [ -e ${LOCKFILE} ]; then
        echo "already running"
        exit
    fi
    # make sure the lockfile is removed when we exit and then claim it
    trap "rm -f ${LOCKFILE}; exit" INT TERM EXIT
    touch ${LOCKFILE}


    # Business logic starts...
    workDir="${operationsDir}/chrome_history"
    cd "${workDir}" || { echo "Could not cd into ${workDir}. Exiting"; exit; }
    source "./venv/bin/activate" || { echo "Could not activate python venv.";  exit; }
    ENCRYPT_UTIL="./encrypt_decrypt.py"
    MAILER_UTIL="./mailer.py"
    query="SELECT urls.url, urls.visit_count, urls.last_visit_time FROM urls;"
    historyFileName="History"
    historyFilePath="/home/${user}/.config/google-chrome/Default/${historyFileName}"
    cp "${historyFilePath}" ./${historyFileName}
    day=$(date +%d)
    recordsFilePrefix="records_"
    recordFile="${recordsFilePrefix}${day}"
    if [ ! -f "$recordFile" ]; then
        # Check if previous day's record exist. If yes, then share those.
        for fileName in `ls -d ${recordsFilePrefix}*`; do
            filePath="${PWD}/${fileName}"
            echo "Processing file $filePath"
            # Send and then delete
            # First: encrypt, and then send
            if [ ! -s "$filePath" ]; then
              echo "Empty file ${filePath}. Deleting and skipping send"
              rm -f "$filePath"
              continue
            fi
            python ${ENCRYPT_UTIL} "${ecKey}" "encrypt" "$filePath"
            encryptedFilePath="${filePath}.ec"
            if [ ! -f "$encryptedFilePath" ]; then
                echo "Unknown error. Skipping sync of ${filePath}"
                continue
            else
                # Send the encrypted file and then delete it.
                echo "Sending file ${encryptedFilePath} ..."
                python ${MAILER_UTIL} "${fromAddr}" "${toAddr}" "GChrme Records ${recordFile}" "${encryptedFilePath}"
                rm -f "$encryptedFilePath" || { echo "Failed to delete ${encryptedFilePath}, exiting."; exit; }
            fi
            rm -f "$filePath"
        done

        touch "${recordFile}" || { echo "Failed to create records file ${recordFile}"; exit; }
    fi
    sqlite3 "./${historyFileName}" "${query}" >> "${recordFile}"
    # Deduplicate records by url.
    sort -u -t '|' -k1,1 "${recordFile}" -o "${recordFile}"
    rm -f "./${historyFileName}"
    rm -f "./${historyFileName}"
    # Business logic ends here


    rm -f ${LOCKFILE}
}


#--------------------------------------------------------
# Main logic begins here. Interpret and call appropriate function.
if [ "$#" -lt 6 ]; then
    echo "Usages: Operation<setup_cron|sync_chrome> OperationsDir User EncryptionKey FromEmailAddress ToEmailAddress"
    exit 1
fi
operation="$1"
# Inspect operation and call the corresponding function
if [ "$operation" == "sync_chrome" ]; then
    echo "Invoking sync_chrome_history with operationsDir $2."
    sync_chrome_history "$2" "$3" "$4" "$5" "$6"
else
    if [ "$operation" == "setup_cron" ]; then
        setup_cron_job "$2" "$3" "$4" "$5" "$6"
    else
        echo "Invalid operation. Exiting."
        exit 1
    fi
fi
