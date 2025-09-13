# poc-api-key

## Description

This is a proof-of-concept (POC) for using an API key to secure internal RESTful services in Go. It demonstrates a simple approach for internal APIs where full OAuth2, JWT, or OpenID Connect may be unnecessary.

## Reason

For an **internal service** (used only inside your company), you don’t always need the _full weight_ of OAuth2, JWTs, or OpenID Connect unless you already use them company-wide. The trade-off is between **simplicity** and **security needs**.

## API Key Simulation

In real-life scenarios, the API key should be set and managed only on the server side. Clients do not send the API key; the server appends it automatically for trusted internal requests.

### Best Practices

- **Store API keys securely**: Use environment variables or a secret manager. Avoid hardcoding keys in source code or configuration files.
- **Rotate API keys periodically**: Change keys on a regular schedule and after any suspected compromise.
- **Restrict network access**: Use firewalls, VPNs, or private networks to limit who can reach your internal service.
- **Monitor usage**: Log and audit API key usage to detect anomalies or unauthorized access.

> **Note:** If your service ever becomes public-facing or needs to interact with external clients, consider upgrading to stronger authentication mechanisms such as OAuth2, JWT, or OpenID Connect.

## Setup

To get started quickly, use:

```sh
make setup
```

This will prepare the project dependencies so you can run:

```sh
go run main.go
```

## REST Client (Optional)

For a better experience testing the API, install the [REST Client extension](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) for VS Code. You can use the provided `.rest` file to interact with the API endpoints easily.
