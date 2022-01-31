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
gotz --timezones "Office:America/New_York,Home:Europe/Berlin"
```

(lookup timezones in the [timezones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) wiki page - _TZ database name_ column)

Set 12-hour format:

```bash
gotz --hours12 true
```

## Customization

The configuration is stored in `$HOME/.gotz.config.json`. It can be configured directly or via the arguments of the `gotz` command (see `gotz --help`). The configuration attributes are described in the following example:

```jsonc
{
    // Configures how the day is segmented
    "day_segments": {
        // Hour of the morning to start (0-23)
        "morning": 6,
        // Color of the morning segment (named color or terminal color code)
        "morning_color": "red",
        // Hour of the day (business hours / main time) to start (0-23)
        "day": 8,
        // Color of the day segment (named color or terminal color code)
        "day_color": "yellow",
        // Hour of the evening to start (0-23)
        "evening": 18,
        // Color of the evening segment (named color or terminal color code)
        "evening_color": "blue",
        // Hour of the night to start (0-23)
        "night": 22,
        // Color of the night segment (named color or terminal color code)
        "night_color": "red"
    },
    // Configures the timezones to be shown
    "timezones": [
        // Timezones have a name (Name) and timezone code (TZ)
        { "Name": "Office", "TZ": "America/New_York" },
        { "Name": "Home", "TZ": "Europe/Berlin" },
    ],
    // Select symbols to use for the time blocks (one of 'mono', 'rectangles' or 'sun-moon')
    "symbols": "rectangles",
    // Indicates whether to plot tics for the local time
    "tics": false,
    // Indicates whether to stretch across the full terminal width (causes inhomogeneous segment lengths)
    "stretch": true,
    // Indicates whether to colorize the blocks
    "colorize": false,
    // Indicates whether to use 12-hour format
    "hours12": false
}
```

## Why?

Working in an international team is a lot of fun, but comes with the challenge of having to deal with timezones. Since I am not good at computing them quickly in my head, I decided to write a simple CLI tool to help me out. I hope it can be useful for other people as well.
Thanks for the inspiration @[sebas](https://github.com/sebastian-quintero)!
