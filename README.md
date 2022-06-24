# nyx

Nyx is a command line tool for checking the status of multiple resource types.

## Installation

Ensure you have `$GOPATH/bin` directory in your `$PATH` to use this tool.
Install the latest version of `nyx` by running:
```shell
$ go install github.com/fsrv-xyz/nyx/cmd/nyx@latest
```

## Features

* checks
  * tcp port reachability for local and remote hosts
  * process status on base of pidfiles
  * exit status of shell commands
  * ssl certificate validity and expiration
  * what ever you want to implement
  
* UI
  * terminal output in colored table format
  * output filtering based on check `identifier`

## Configuration

Nyx is configurable via a json file. The default configuration file is `./nyx.json` and can be adjusted by setting the `NYX_CONFIG` environment variable or `-config.file` parameter.