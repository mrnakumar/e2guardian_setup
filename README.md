### How to install
The file setup.sh can be executed to install e2guardian-5.4.
Run using the command `./setup.sh`

### All python related code can be run by doing the folllowing:
* Create a python3 virtual environment.
* Install required packages by running `pip install -r ./requirements-for-python-code.txt`

### Using encrypt_decrypt
#### To encrypt a file named `iptables.service` in current working directory:
```shell
python encrypt_decrypt.py "mysimplekey" "encrypt" "./iptables.service"
```

#### To decrypt a file named `iptables.service.ec` in current working directory:
```shell
python encrypt_decrypt.py "mysimplekey" "decrypt" "./iptables.service.ec"
```

### Using mailer
```shell
python mailer.py "fromEmailAddress" "toEmailAddress" "emailSubject" "filePathToSendAsAttachment"
```

#### Install chrome-history-sync script as cron job
Use the below given command to install the cron job for a specific user:
```shell
./cron_chrome_history_sync.sh "setup_cron" "infinity"
```
Note:- Before executing the above command make sure that:
* The given user (i.e. ${user} exist on that computer
* The directory `/etc/${user}/chrome_history` must exist and ${user} should have write permission on that
* The current working directory must have the following:
mailer.py
encrypt_decrypt.py
chrome-history-sync.sh
requirements-for-python-code.txt
token.json (this is the gmail API Oauth token)
