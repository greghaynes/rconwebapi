rconwebapi
==========

Simple HTTP/REST JSON bridge for CS:GOs [RCON protocol](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol).


# Usage

Releases can be found [here](https://github.com/greghaynes/rconwebapi/releases).

Run the server binary:

```bash
./rconwebapi -host 127.0.0.1:8099
```


# API

## POST

POST to `/rcon` with the following Body:

```json
{
  "RconRequest": {
    "Address": "myserver.com:port",
    "Password": "secret_password",
    "Command": "status"
  }
}
```

Get a response with the folowing Body:

```json
{
    "RconResponse": {
        "Output":"hostname: ...\nversion : 1.37.5.2/13752....\n"
    }
}
```

## Websocket

Open a websocket on `/rcon_ws` and send the following text messages:

### Connect

Request

```json
{
    "RequestType": "connect",
    "Request": {
        "Address": "myserver.com:port",
        "Password": "secret_password"
    }
}
```

Response

```json
{
    "ResponseType": "connect",
    "Response": {
        "Status": "success/fail",
        "Message": "maybe something"
    }
}
```

### Command

Request

```json
{
    "RequestType": "command",
    "Request": {
        "Command": "status"
    }
}
```

Response

```json
{
    "ResponseType": "command",
    "Response": {
        "Output": "something"
    }
}
```