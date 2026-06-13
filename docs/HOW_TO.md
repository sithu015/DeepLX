# DeepLX How-To Guides

This document provides problem-oriented recipes for common deployment, configuration, and security scenarios when operating DeepLX.

---

## How to Secure DeepLX with an Access Token

By default, anyone who can reach your DeepLX endpoint can make translation requests. If you are exposing your instance to the public internet, you should secure it using an access token.

### 1. Set the Token Environment Variable
Start the server with the `TOKEN` environment variable set:
```bash
export TOKEN="my-super-secret-passphrase"
go run main.go
```
Or start via Docker:
```bash
docker run -d -p 1188:1188 -e TOKEN="my-super-secret-passphrase" --name deeplx-app deeplx-local:latest
```

### 2. Make Authenticated Requests
Subsequent calls to `/translate`, `/v1/translate`, and `/v2/translate` will fail with `401 Unauthorized` unless you supply the correct token.
Provide the token in the `Authorization` header as a Bearer token:
```bash
curl -i -X POST http://localhost:1188/translate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer my-super-secret-passphrase" \
  -d '{
    "text": "Hello, world!",
    "target_lang": "DE"
  }'
```

---

## How to Configure Pro-Tier Translation with dl-session

If you have a paid DeepL Pro subscription, you can configure DeepLX to run translation queries through the DeepL Pro oneshot API, which has no character limitations.

### 1. Retrieve your OAuth Access Token
1. Log in to your paid DeepL Pro account on a web browser.
2. Open your browser's Developer Tools (F12) -> Application/Storage tab.
3. Locate the Cookie values under `www.deepl.com` and find the value of the `dl_access` cookie (an OAuth JWT token starting with `eyJ...`).

### 2. Start DeepLX with Pro Credentials
Configure the server with the `DL_SESSION` environment variable or the `-s` flag:
```bash
export DL_SESSION="eyJhbGciOi..."
go run main.go
```
Or run with Docker:
```bash
docker run -d -p 1188:1188 -e DL_SESSION="eyJhbGciOi..." --name deeplx-app deeplx-local:latest
```

### 3. Make Translation Requests
Send requests to `/v1/translate`. DeepLX will route these requests to `oneshot-pro.www.deepl.com` using your token.
```bash
curl -i -X POST http://localhost:1188/v1/translate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "This is a very long text that will bypass the standard free-tier length limitations...",
    "target_lang": "ES"
  }'
```

---

## How to Decrypt and Extract the dl_access Session Cookie on macOS

On macOS, Google Chrome encrypts its SQLite cookie database using a key stored in the macOS Keychain. You can extract the `dl_access` token programmatically using the following steps:

### 1. Retrieve the Chrome Safe Storage Password
Run the following command in your terminal to extract the raw encryption password from the Keychain:
```bash
security find-generic-password -w -s "Chrome Safe Storage"
```

### 2. Decrypt the Cookie with Python
Use the password (e.g. `YBeMNd2hq3/kAd1F6+bNrg==`) in the following Python script to copy, query, and decrypt the `dl_access` token:
```python
import sqlite3
import os
import hashlib
import subprocess
import binascii

# Derived from keychain query
password = b"YBeMNd2hq3/kAd1F6+bNrg=="
salt = b"saltysalt"
iterations = 1003
key = hashlib.pbkdf2_hmac('sha1', password, salt, iterations, dklen=16)

def decrypt_openssl(ciphertext, key_bytes):
    iv_hex = "20202020202020202020202020202020"
    key_hex = binascii.hexlify(key_bytes).decode()
    proc = subprocess.Popen(
        ['openssl', 'enc', '-d', '-aes-128-cbc', '-K', key_hex, '-iv', iv_hex, '-nosalt'],
        stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE
    )
    stdout, _ = proc.communicate(input=ciphertext)
    if proc.returncode == 0:
        decrypted = stdout
        if len(decrypted) > 0:
            padding_len = decrypted[-1]
            if padding_len < 16:
                decrypted = decrypted[:-padding_len]
        return decrypted
    return b""

# Copy to avoid SQLite lock
cookies_path = os.path.expanduser("~/Library/Application Support/Google/Chrome/Default/Cookies")
temp_db_path = "/tmp/cookies_temp.db"
import shutil
shutil.copyfile(cookies_path, temp_db_path)

conn = sqlite3.connect(temp_db_path)
cursor = conn.cursor()
cursor.execute("SELECT name, encrypted_value FROM cookies WHERE name='dl_access' LIMIT 1")
row = cursor.fetchone()
if row:
    name, enc_val = row
    decrypted_bytes = decrypt_openssl(enc_val[3:], key)
    # The first 32 bytes are a signature/header; the rest is the plain-text JWT token.
    token = decrypted_bytes[32:].decode('utf-8', errors='ignore')
    print("dl_access token:", token)

conn.close()
os.remove(temp_db_path)
```
You can save this decrypted token directly into your Dokploy environment variables under `DL_SESSION` to enable Pro features.

---


---

## How to Set Up an Outbound Proxy

DeepL aggressively rate-limits IP addresses that make excessive translation requests on the free tier. If your server is blocked, you can route outbound traffic through an HTTP/HTTPS proxy.

### Start DeepLX with Outbound Proxy
Specify the proxy URL using the `PROXY` environment variable or `-proxy` CLI flag:
```bash
export PROXY="http://user:password@proxy-server.com:8080"
go run main.go
```
Or run with Docker:
```bash
docker run -d -p 1188:1188 -e PROXY="http://user:password@proxy-server.com:8080" --name deeplx-app deeplx-local:latest
```
DeepLX will parse this proxy configuration and allocate an isolated HTTP client pool, routing all requests to `oneshot-free.www.deepl.com` through the proxy.

---

## How to Install DeepLX as a systemd Service

To ensure DeepLX automatically runs on boot and restarts on crash, you can deploy it as a systemd service.

### 1. Download the Installation Scripts
```bash
wget https://raw.githubusercontent.com/sithu015/DeepLX/main/install.sh
wget https://raw.githubusercontent.com/sithu015/DeepLX/main/uninstall.sh
chmod +x install.sh uninstall.sh
```

### 2. Run the Installer as Root
```bash
sudo ./install.sh
```
The script will download the latest binary matching your architecture, install it to `/usr/bin/deeplx`, configure a dedicated, unprivileged system user (`deeplx`), and register the service.

### 3. Control the Service
```bash
# Check status
sudo systemctl status deeplx

# Restart
sudo systemctl restart deeplx

# Edit environment variables (e.g. set TOKEN/PROXY)
sudo systemctl edit deeplx
```
Inside the override configuration file, you can specify environment variables:
```ini
[Service]
Environment="TOKEN=my_secret_token"
Environment="PROXY=http://127.0.0.1:8080"
```
Save, close, and reload the service:
```bash
sudo systemctl daemon-reload
sudo systemctl restart deeplx
```
