# duckhunt
Block HID/Keystroke Injection/Rubber Ducky Attacks by measuring the typing speed and blocking all inputs if a threshold is exceeded.

# Installation
## Binaries
Download the pre built binaries [here](https://github.com/rtfmkiesel/duckhunt/releases).
## With go
```
go install github.com/rtfmkiesel/duckhunt@latest
```

# Build from source
```bash
git clone https://github.com/rtfmkiesel/duckhunt
cd duckhunt
# duckhunt.exe (normal, hidden window)
go build -ldflags="-s -w -H=windowsgui" -o "duckhunt.exe" cmd/duckhunt/duckhunt.go
# duckhunt_cli.exe (shows output in terminal)
go build -ldflags="-s -w" -o "duckhunt_cli.exe" cmd/duckhunt/duckhunt.go
```

# Usage
```
duckhunt.exe [OPTIONS]

Options:
 -r Registers duckhunt as a scheduled task to start upon login
```

## Config
```yaml
# duckhunt config file
# alert via popup
alert: true
# alert in foreground
alertontop: true
# popup title
alerttitle: "WARNING"
# popup message
alertmessage: "Fast keystrokes detected, check your USB Ports"

# block user inputs if fast keystrokes are detected, duration of block in seconds
# to not block set value to 0
blockduration: 10

# name for the logfile, will get saved into the same folder as the executable
logfile: "duckhunt.log"

# ADVANCED
# maximal time interval between keystrokes in milliseconds
maxInterval: 50

# how many keystrokes will be used to calculate the average time interval
# this value sets the minimum length of a payload. if an attacker has a payload with less keystrokes than this value, it will succeed
historySize: 25

# some ignored keys like: SPACE, SLASH, CR, DEL, backspace, arrow keys, media keys etc.
# key that are often sent repeatedly
# see https://www.toptal.com/developers/keycode/table
ignoredKeys: [8, 13, 32, 47, 90, 91, 92, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183]`
```

# Kudos
The concept is taken from [pmsosa/duckhunt](https://github.com/pmsosa/duckhunt).