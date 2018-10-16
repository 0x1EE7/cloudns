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

```shellsession
$ cloudns add -i `curl -s https://ipinfo.io/ip`  -d your.app.tld
Using config file: ~/.cloudns.yaml
Adding IPs [2.2.3.4] to your.app.tld
Up to date records after changes: [2.2.3.4 1.2.3.4]

$ cloudns remove -i 1.2.3.4 -i 2.2.3.4 -d your.app.tld
Using config file: ~/.cloudns.yaml
Removing IPs [1.2.3.4 2.2.3.4] from your.app.tld
Up to date records after changes: []
```

### Config File & env
`sa_file` can be configured in `--config` file
```yaml
$ cat ~/.cloudns.yaml
sa_file: /path/to/google_cloud_sa_file.json
dns_zone: yourapp-zone-name
```
Can also be overriden by an env var
```shellsession
$ SA_FILE=/path/to/anothet_sa.json cloudns
```


### Help
```shellsession
$ cloudns
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

