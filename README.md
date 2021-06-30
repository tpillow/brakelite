# Brakelite
A lightweight, minimal, and effective notification system to remind you to stretch
and take quick breaks throughout the day.

## How it Works
Brakelite opens as a system tray icon, and all options can be found in the system
tray drop-down menu.

### Features
- Pre-set notification durations to choose from
- Pre-set notification pause durations to choose from
- Next notification/pause status indicator
- Hard-coded, randomized break messages

### Future Plans
- Unit tests (will need to mock notifier and add time scale)
- Automatic build/release system for executables (Github)

### Building Executables
Brakelite can be built for any Go OS/ARCH, for example:
- env GOOS=windows GOARCH=amd64 go build build/win/amd64/brakelite.exe
- env GOOS=linux GOARCH=amd64 go build build/linux/amd64/brakelite
