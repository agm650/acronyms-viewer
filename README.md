# acronyms-viewer

![Go Test](https://github.com/agm650/acronyms-viewer/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/agm650/acronyms-viewer/workflows/goreleaser/badge.svg)

## Table of Contents

- [acronyms-viewer](#acronyms-viewer)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Usage examples](#usage-examples)
  - [Configuration](#configuration)
  - [Installation from source](#installation-from-source)

## Overview

This is a basic web app used to search for acronyms definitions.
Acronyms are distributed in the for of a SQLite database that the application will use.

## Usage examples

In order to use this tool you will have to provide many parameters:

```bash
#  /usr/local/bin/acronyms-viewer -h
TODO

Usage:
  acronyms-viewer [flags]

Flags:
      --config string           config file (default is $HOME/.acronyms-viewer.yaml)
  -d, --debug                   Activate debug messages
  -h, --help                    help for acronyms-viewer
  -D, --database                path to the database to use
  -P, --port                    port to listen for HTTP request (default: 8080)
```

Here is an example of the command line to use:

```bash
/usr/local/bin/acronyms-viewer -D ~/acro_database.db -P 9090
```

## Configuration

All parameters can be provided as flags on the command line, or as environment variable, or in a configuration file.

Configuration file example:

```yaml
---
debug: false
database: ~/acro_database.db
port: 9090
```

## Installation from source

From the local path of the acronyms-viewer repository:

```bash
go build -o bin/acronyms-viewer
sudo cp bin/acronyms-viewer /usr/local/bin/
sudo chmod 555 /usr/local/bin/acronyms-viewer
```

Then copy the resulting executable at the desired location on your system.
