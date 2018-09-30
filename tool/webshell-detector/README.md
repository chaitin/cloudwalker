# Webshell detector

## Installation

Standalone binary can be downloaded at <https://github.com/chaitin/cloudwalker/releases>.

## Usage

```
./webshell-detector-linux-amd64 -h
Chaitin CloudWalker Webshell Detector
[version 1.0.0]

usage: ./webshell-detector-linux-amd64 [options] name ...

  -html
    	Show result as HTML
  -output string
    	Export result to output file
```

## Build

### Dependencies

For Ubuntu and Debian users:

```
apt-get install autoconf bzip2 patch vim
```

### Build the detector

```
make -C php
go get -d ./src
go build -o webshell-detector ./bin
```
