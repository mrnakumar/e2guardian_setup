### How to install
The file setup.sh can be executed to install e2guardian-5.4.
Run using the command `./setup.sh`

### To use mailer create python virutal environment, activate and install the following python packages:
pip install --upgrade google-api-python-client google-auth-httplib2 google-auth-oauthlib

### To use encrypt_decrupt install the following python packages:
pip install pycryptodome

### Using encrypter_decrypter
#### To encrypt a file named `iptables.service` in current working directory:
```shell
python encrypt_decrypt.py "mysimplekey" "encrypt" "./iptables.service"
```

#### To decrypt a file named `iptables.service.ec` in current working directory:
```shell
python encrypt_decrypt.py "mysimplekey" "decrypt" "./iptables.service.ec"
```
