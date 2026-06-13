# DeepLX Configuration & API Reference

This document provides a formal reference for all CLI flags, environment variables, API endpoints, response schemas, and supported languages in DeepLX.

---

## 1. CLI Configurations

Configurations are set up via environment variables or CLI flags. If both are specified, the CLI flag takes precedence.

| Environment | CLI Flag | Default | Description |
| :--- | :--- | :--- | :--- |
| `IP` | `-ip`, `-i` | `0.0.0.0` | Bind IP address. |
| `PORT` | `-port`, `-p` | `1188` | Listen port. Value must be a valid integer in `[1, 65535]`. |
| `TOKEN` | `-token` | `""` | Optional authorization passphrase required by callers. |
| `DL_SESSION` | `-s` | `""` | DeepL Pro OAuth access token (JWT string) used for `/v1/translate`. |
| `PROXY` | `-proxy` | `""` | Outbound proxy URL (e.g. `http://127.0.0.1:1080` or `socks5://...`). |

---

## 2. API Endpoints

All endpoints respond with standardized JSON bodies. If authentication fails, the server responds with a `401 Unauthorized` status code.

### 2.1. GET `/`
Root endpoint returning project details.
- **Request Headers**: None.
- **Response Payload**: `200 OK`
  ```json
  {
    "code": 200,
    "message": "DeepL Free API, Developed by sjlleo and missuo. Go to /translate with POST. http://github.com/OwO-Network/DeepLX"
  }
  ```

---

### 2.2. POST `/translate`
Main translation endpoint on the free tier.
- **Text Length Limit**: Total length of the `text` array/input must not exceed 1500 characters.
- **Request Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer <TOKEN>` (Required only if `TOKEN` configuration is non-empty)
- **Request Body JSON Schema**:
  ```json
  {
    "text": "String, required. Text to translate.",
    "source_lang": "String, optional. Defaults to \"auto\".",
    "target_lang": "String, required. Target BCP-47 language code.",
    "tag_handling": "String, optional. Allowed values: \"html\", \"xml\"."
  }
  ```
- **Response Payload JSON Schema**: `200 OK`
  ```json
  {
    "code": 200,
    "id": 1718251293000,
    "data": "Translated text string.",
    "alternatives": null,
    "source_lang": "Detected source language code (e.g. \"EN\").",
    "target_lang": "Target language code (e.g. \"DE\").",
    "method": "Free"
  }
  ```
- **Common Error Codes**:
  - `400 Bad Request`: `Invalid request payload` or invalid parameter constraints.
  - `401 Unauthorized`: `Invalid access token` (Authorization header mismatch).
  - `413 Payload Too Large`: `text exceeds maximum length` (Input exceeds 1500 limit).
  - `429 Too Many Requests`: `too many requests, your IP has been blocked by DeepL temporarily`.

---

### 2.3. POST `/v1/translate`
Translation endpoint requiring Pro-tier sessions.
- **Text Length Limit**: None (Bypasses character length checks).
- **Request Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer <TOKEN>` (Required only if `TOKEN` configuration is non-empty)
  - `Cookie: dl_session=<PRO_OAUTH_TOKEN>` (Session cookie, takes precedence over the configured `DL_SESSION`)
- **Request Body**: Same as `/translate`.
- **Response Payload**: Same as `/translate` with `method: "Pro"`.

---

### 2.4. POST `/v2/translate`
Officially-compatible DeepL API endpoint.
- **Request Headers**:
  - `Content-Type: application/json` or `application/x-www-form-urlencoded`
  - `Authorization: Bearer <TOKEN>` (Required only if `TOKEN` is set)
- **Request Body (JSON or Form Fields)**:
  ```json
  {
    "text": ["Array of strings to translate. Required."],
    "target_lang": "String, target BCP-47 language code. Required."
  }
  ```
- **Response Payload**: `200 OK`
  ```json
  {
    "translations": [
      {
        "detected_source_language": "Detected language code (e.g. \"EN\").",
        "text": "Translated text string."
      }
    ]
  }
  ```

---

## 3. Supported Languages

DeepLX supports all standard source and target languages made available by DeepL. If an unknown BCP-47 code is passed, it is forwarded directly to the upstream server.

### Supported Source Languages (`source_lang`)
- `AR` (Arabic)
- `BG` (Bulgarian)
- `CS` (Czech)
- `DA` (Danish)
- `DE` (German)
- `EL` (Greek)
- `EN` (English - regional variant autodetected)
- `ES` (Spanish)
- `ET` (Estonian)
- `FI` (Finnish)
- `FR` (French)
- `HE` (Hebrew)
- `HU` (Hungarian)
- `ID` (Indonesian)
- `IT` (Italian)
- `JA` (Japanese)
- `KO` (Korean)
- `LT` (Lithuanian)
- `LV` (Latvian)
- `MY` (Burmese)
- `NB` (Norwegian Bokmål)
- `NL` (Dutch)
- `PL` (Polish)
- `PT` (Portuguese - regional variant autodetected)
- `RO` (Romanian)
- `RU` (Russian)
- `SK` (Slovak)
- `SL` (Slovenian)
- `SV` (Swedish)
- `TH` (Thai)
- `TR` (Turkish)
- `UK` (Ukrainian)
- `VI` (Vietnamese)
- `ZH` (Chinese - simplified)

### Supported Target Languages (`target_lang`)
Regional variants must be specified when translating *into* English or Portuguese:
- `EN-US` (English - American)
- `EN-GB` (English - British)
- `PT-BR` (Portuguese - Brazilian)
- `PT-PT` (Portuguese - European)
- `ZH-HANS` (Chinese - Simplified)
- `ZH-HANT` (Chinese - Traditional)
- All other base codes listed under source languages (e.g. `DE`, `FR`, `JA`, etc.) are supported directly.
