# Simple Opt Out Example

This example shows how to setup simple opt-out verification.

## Setup

### 1. Download Repo

```
git clone https://github.com/pBread/Simple-Opt-Out.git
cd simple-opt-out
```

### 2. Purchase a phone number in your Twilio Console.

### 3. Add the environment variables per [.env.example](./.env.example)

```
ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
AUTH_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
OPT_OUT_PHONE=+12223334444
```

### 4. Start Go

```
go run .
```

### 5. Create ngrok tunnel on :8080 (or create a public URL in some other way)

```
ngrok http 8080
```

<img src="https://i.imgur.com/TU1vJSf.png"/>

### 6. Assign the inbound SMS webook to the `/opt-out/reply` route

<img src="https://i.imgur.com/LtTAPFt.png" height="400" />

## Usage

Initiate the verification process by going to `http://localhost:8080/opt-out/new?to=+12223334444`, but replace the phone number with your personal phone number. You should receive an SMS message.

Reply to the text message with and w/out the keyword `TESTER` to see what happens.

NOTE: This example does not unsubscribe you when you respond with the keyword. You need to update [func unsubscribe](./main.go#L81) with the necessary logic to disable the user/post in your database.
