# Meta Cloud API Setup Guide

## Prerequisites

* Meta account
* Business Manager account
* WhatsApp Business number

## Step 1: Create App

1. Access Meta Developers
2. Create Business App
3. Add WhatsApp product

## Step 2: Configure Webhook

Endpoint:

```http
POST /webhooks/whatsapp
```

Verification endpoint:

```http
GET /webhooks/whatsapp
```

Validation params:

* hub.mode
* hub.challenge
* hub.verify_token

## Step 3: Generate Access Token

Store token securely.

## Step 4: Capture IDs

Required values:

* phone_number_id
* business_account_id

## Environment Variables

```env
WHATSAPP_ACCESS_TOKEN=
WHATSAPP_PHONE_NUMBER_ID=
WHATSAPP_VERIFY_TOKEN=
META_APP_SECRET=
```

## Send Message Example

```http
POST https://graph.facebook.com/v22.0/{phone-number-id}/messages
```

Payload:

```json
{
  "messaging_product": "whatsapp",
  "to": "559199999999",
  "type": "text",
  "text": {
    "body": "Olá"
  }
}
```

