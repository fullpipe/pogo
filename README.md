# pogo

Pogo is a simple [pomodoro
technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) tool.
You start new session and nothing will distract you.

## Installation

```bash
go get -u github.com/fullpipe/pogo
```

## Usage

Everything that pogo does is modifying `/etc/hosts`. So it requires write
permissions. I persоnaly start it using sudo.

```bash
sudo pogo
```

or to run it in background

```bash
nohup sudo pogo &
```

then start your pomodoro session from tray

![systray example](tray.png "systray example")

## Configuration

By default pogo's session contains 4 pomodoro sessions: 25 mins of work and 5
mins of rest.  
This could be modified by config file.

```yaml
# ~/.config/pogo/pogo.yaml
pomodoro:
  work: 25 # work time
  rest: 5 # rest time
  repeats: 4
work_ip: 127.0.0.1 # ip address for distracting domains
domains:
  good: # list of good domains that should be always available
    - youtube.com
  bad: # list of bad domains which will be added to basic distracting domains
    - facebook.com
```

[Basic distracting
domains](https://github.com/fullpipe/pogo/blob/master/config.go#L103)

## TODO

- wrap it in app?
- no sudo solution
- icons
