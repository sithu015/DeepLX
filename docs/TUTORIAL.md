# DeepLX Tutorial: Quick Start Guide

This tutorial guides you through setting up and running your own DeepLX instance, testing it locally, and making your first translation request.

---

## Prerequisites
To follow this tutorial, you need one of the following installed:
- **Option A (Native)**: [Go](https://go.dev/doc/install) (version 1.25 or later) installed on your system.
- **Option B (Docker)**: [Docker](https://docs.docker.com/get-docker/) installed and running.

---

## Step 1: Clone the Repository
Open a terminal and run the following command to download the source code:
```bash
git clone https://github.com/sithu015/DeepLX.git
cd DeepLX
```

---

## Step 2: Choose Your Running Method

### Option A: Running Locally (Native Go)
Build and run the executable using the standard Go toolchain:
```bash
# Compile and run in one step
go run main.go
```
By default, the server starts up and begins listening on `0.0.0.0:1188`. You will see output similar to this:
```text
DeepL X has been successfully launched! Listening on 0.0.0.0:1188
Developed by sjlleo <i@leo.moe> and missuo <me@missuo.me>.
```

---

### Option B: Running via Docker
If you prefer Docker, you can compile and build the container locally:
```bash
# Build the Docker image
docker build -t deeplx-local .

# Run the container and map port 1188
docker run -d -p 1188:1188 --name deeplx-app deeplx-local:latest
```

---

## Step 3: Verify the Server is Running
Test that the server is alive by making a simple HTTP GET request to the root endpoint `/`:
```bash
curl -i http://localhost:1188/
```
You should receive a `200 OK` response with a JSON payload:
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{
  "code": 200,
  "message": "DeepL Free API, Developed by sjlleo and missuo. Go to /translate with POST. http://github.com/OwO-Network/DeepLX"
}
```

---

## Step 4: Make Your First Translation
DeepLX exposes a free translation API via `POST /translate`. Let's translate "Hello, world!" from English (`EN`) to German (`DE`):

```bash
curl -i -X POST http://localhost:1188/translate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, world!",
    "source_lang": "EN",
    "target_lang": "DE"
  }'
```

You should receive a translation response containing the translated text under `data`:
```json
{
  "code": 200,
  "id": 1718251293000,
  "data": "Hallo, Welt!",
  "alternatives": null,
  "source_lang": "EN",
  "target_lang": "DE",
  "method": "Free"
}
```

---

## Next Steps
Congratulations! You have set up and verified your own translation proxy server.
- To secure your endpoint with an access token, see the [How-To Guide on Security](HOW_TO.md).
- For complete API documentation, consult the [Reference Guide](REFERENCE.md).
- To understand how the client connection pooling and WAF bypassing works behind the scenes, read the [Explanation of System Design](EXPLANATION.md).
