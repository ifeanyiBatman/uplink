# Uplink

Uplink is a lightweight, simplified CLI tool built on top of the official `ngrok-go` library. It allows you to easily bypass NATs and firewalls to expose your local web servers to the internet using custom ngrok domains—without needing to install or manage the full ngrok daemon.

## Features

- **Quick & Easy Tunneling**: Instantly forward your local development servers to the public internet.
- **Persistent Configuration**: Saves your credentials (authtoken, custom domain, and username) locally to `~/.uplinkconfig.json` so you only have to configure it once.
- **Custom Ports**: Easily specify which localhost port to expose using the `-port` flag (defaults to `8080`).
- **Clean Shutdown**: Type `quit` into the terminal at any time for a graceful shutdown of the tunnel.

## Prerequisites

- Go 1.26 or higher
- An [ngrok](https://ngrok.com/) account
- An ngrok **Authtoken** (get it from the [ngrok dashboard](https://dashboard.ngrok.com/get-started/your-authtoken))
- A claimed **static domain** on ngrok (e.g., `your-custom-name.ngrok.app`)

## Installation

You can clone the repository and build the binary:

```bash
git clone https://github.com/ifeanyiBatman/uplink.git
cd uplink
go build -o uplink main.go
```

Alternatively, you can install it directly via `go install` (assuming your Go workspace `/bin` is in your `PATH`):

```bash
go install github.com/ifeanyiBatman/uplink@latest
```

## First Run Configuration

The first time you run Uplink, you will be prompted to enter your ngrok credentials. These will be securely saved into `.uplinkconfig.json` in your home directory.

```text
$ ./uplink
Please setup Ngrok... https://dashboard.ngrok.com/get-started/setup
Enter your authtoken: <your_ngrok_authtoken_here>
Enter your domain: <your-custom.ngrok.app>
Enter your username: <your_preferred_username>
```

## Usage

Once configured, simply run `uplink` and specify the port your local server is running on.

Expose the default port (8080):
```bash
./uplink
```

Expose a custom port (e.g., 3000):
```bash
./uplink -port 3000
```

### Stopping the Tunnel

While the tunnel is running, you will see a prompt:

```text
Type quit to shutdown
Tunnel is running. Type 'quit' to shutdown.
Your machine here http://localhost:3000 has been forwarded to your-custom.ngrok.app
```

Simply type `quit` and hit Enter to gracefully close the tunnel.
