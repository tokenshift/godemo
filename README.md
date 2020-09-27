# Godemo

Command-line tool for running command-line demos.

## Installation

```
$ go get -u github.com/tokenshift/godemo
```

## Usage

```
$ godemo --help
$ godemo setup -f demo.yaml
$ godemo run -f demo.yaml
$ godemo teardown -f demo.yaml
```

The `godemo` tool runs demo steps, broken into `setup`, `run`, and `teardown`
phases. Each command will run the demo steps for that phase, prompting for input
after each step. By default, simply hitting 'enter' will continue to the next
step, but you can override this behavior by entering one of the following options:

* `(n)ext`   - Proceed to the next step (default)
* `(p)prev`  - Go back to the last step
* `(r)eplay` - Repeat the current step
* `(q)uit`   - Exit the demo
* `{number}` - Goto a specific step

You will also be prompted for input between echoing the command to be run and
actually executing it, so that you can choose to skip that command if you want.

## Demo Definition

```yaml
---
title: Title of the demo
variables:
  # Variables can be used in commands via variable substition, e.g.:
  # curl ${BASE_URL}/index.html
  - name: BASE_URL
    # Name of the variable.
    value: http://localhost:3000
    # Initial value of the variable. Can be overriden later by using `capture`
    # in demo steps (see below).
    echo: true
    # Whether to display the actual value when commands are echoed before
    # running, or just the substitution, i.e. "${BASE_URL}". The default is `true`.
setup:
  # List of preparatory steps that you'll likely run ahead of actually starting
  # your demo. Same structure as `steps`, below.
steps:
  - comment: Display some optional explanatory text.
    # (Optional) text that will be displayed before the command.
    cmd: 'curl ${BASE_URL}/api/example -H "Content-Type: application/json"'
    # (Optional) command to run. Both the comment and command will be displayed
    # before the command is run, with an option to skip or go to an entirely
    # different step.
    echo: false
    # Whether to display the results of this command. Defaults to `true`. Most
    # likely used with...
    capture: NEW_VAR_NAME
    # Capture the output of this command in a new variable, which can be used in
    # subsequent commands. This variable will inherit the commands `echo` setting.
teardown:
  # List of cleanup steps that you'll likely run after the demo is over. Same
  # structure as `steps`, below.
```

The different demo phases are provided as a convenience and are all optional;
they all run the same way, so you could put all of your steps under `steps` and
leave `setup` and `teardown` out.