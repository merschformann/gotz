<img src="material/icon/world.svg" align="right" height="110"/>

# go**tz**

A simple CLI timezone info tool.

## Installation

### Directly via [Go](https://go.dev/doc/install)

```bash
go install github.com/merschformann/gotz@latest
```

### Binary

Simply download the binary of the [latest release](https://github.com/merschformann/gotz/releases/latest/) (look for `gotz_OS_ARCH`), rename it to `gotz` and put it in a folder in your `$PATH`.

## Usage

Show current time:

```bash
gotz
```

![preview](material/screenshot/preview1.png)

Show arbitrary time:

```bash
gotz 15
```

![preview](material/screenshot/preview2.png)

Time can be one of the following formats:

```txt
15
15:04
15:04:05
3:04pm
3:04:05pm
3pm
1504
150405
2006-01-02T15:04:05
```

## Basic configuration

Set the timezones to be used by default:

```bash
gotz --timezones America/New_York,Europe/Berlin
```

## Customization

TODO: Describe configuration options.

## Why?

Working in an international team is a lot of fun, but comes with the challenge of having to deal with timezones. Since I am not good at computing them quickly in my head, I decided to write a simple CLI tool to help me out. I hope it can be useful for other people as well.
Thanks for the inspiration @[sebas](https://github.com/sebastian-quintero)!
