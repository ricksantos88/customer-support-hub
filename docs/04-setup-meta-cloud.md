# Meta Cloud API Setup

## Steps

### 1. Create Meta App
- Access Meta Developers
- Create business app

### 2. Configure WhatsApp Product
- Add WhatsApp product

### 3. Configure Webhook
Webhook URL:

```text
POST /webhooks/whatsapp
```

Verify token required.

### 4. Generate Access Token
Store securely.

### 5. Save IDs
Required:
- phone_number_id
- business_account_id

## Environment Variables

```env
WHATSAPP_TOKEN=
WHATSAPP_PHONE_NUMBER_ID=
WHATSAPP_VERIFY_TOKEN=
```
