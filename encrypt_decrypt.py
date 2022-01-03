import base64
import os
import sys
from Crypto.Cipher import AES
from Crypto.Hash import SHA256
from Crypto import Random


def encrypt(key, source, encode=True):
    key = SHA256.new(
        key
    ).digest()  # use SHA-256 over our key to get a proper-sized AES key
    iv = Random.new().read(AES.block_size)  # generate IV
    encryptor = AES.new(key, AES.MODE_CBC, iv)
    padding = AES.block_size - len(source) % AES.block_size  # calculate needed padding
    source += bytes([padding]) * padding  # Python 2.x: source += chr(padding) * padding
    encrypted = iv + encryptor.encrypt(
        source
    )  # store the IV at the beginning and encrypt
    return base64.b64encode(encrypted).decode("latin-1") if encode else encrypted


def decrypt(key, source, decode=True):
    if decode:
        source = base64.b64decode(source.encode("latin-1"))
    key = SHA256.new(
        key
    ).digest()  # use SHA-256 over our key to get a proper-sized AES key
    iv = source[: AES.block_size]  # extract the IV from the beginning
    decrypter = AES.new(key, AES.MODE_CBC, iv)
    decrypted = decrypter.decrypt(source[AES.block_size :])  # decrypt
    padding = decrypted[
        -1
    ]  # pick the padding value from the end; Python 2.x: ord(data[-1])
    if (
        decrypted[-padding:] != bytes([padding]) * padding
    ):  # Python 2.x: chr(padding) * padding
        raise ValueError("Invalid padding...")
    return decrypted[:-padding]  # remove the padding


def write_binary(contents, filepath):
    g = open(filepath, "wb")
    try:
        g.write(contents)
    finally:
        g.close()


if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usages: encryptionKey operatiion[encrypt|decrypt] sourceFilePath")
        sys.exit(1)

    encryption_key = sys.argv[1]
    operation = sys.argv[2]
    source_file_path = sys.argv[3]
    is_encrypt = True
    folder = os.path.dirname(source_file_path)
    source_filename = os.path.basename(source_file_path)
    target_file_path = ""
    if operation == "decrypt":
        is_encrypt = False
        target_file_path = folder + "/" + source_filename + ".dec"
    else:
        target_file_path = folder + "/" + source_filename + ".ec"

    with open(source_file_path, "rb") as f:
        source_contents = f.read()
        data = None
        if is_encrypt:
            data = encrypt(encryption_key.encode("utf-8"), source_contents, False)
        else:
            data = decrypt(encryption_key.encode("utf-8"), source_contents, False)
        if data is not None:
            write_binary(data, target_file_path)
