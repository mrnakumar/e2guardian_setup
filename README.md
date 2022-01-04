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
