# cloudns

Add & remove dns records to google cloud dns

[![GoDoc](https://godoc.org/github.com/0x1EE7/cloudns/acme?status.svg)](https://godoc.org/github.com/0x1EE7/cloudns)

## Installation

cloudns supports both binary installs and install from source.

To get the binary just download the latest release for your OS/Arch from [the release page](https://github.com/0x1EE7/cloudns/releases)
and put the binary somewhere convenient. cloudns does not assume anything about the location you run it from.

To install from source, just run:

```bash
go get -u github.com/0x1EE7/cloudns
```

## Features

- Add given IPs to the domain
- Remove given IPs from the domain

## Usage

```text
Easily modify DNS records in Google Cloud DNS

cloudns is a CLI to add and remove DNS entries.

Usage:
  cloudns [command]

Available Commands:
  add         Add given IPs to the domain
  help        Help about any command
  remove      Remove given IPs for the domain

Flags:
      --config string   config file (default is $HOME/.cloudns.yaml)
  -h, --help            help for cloudns

Use "cloudns [command] --help" for more information about a command.
```

```text
$ cloudns add -i `curl -s https://ipinfo.io/ip` -d for.bar.tld
Adding IPs [8.8.8.8] to for.bar.tld
```
