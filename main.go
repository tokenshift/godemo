package main

import (
	"github.com/alecthomas/kong"
)

/* Example demo definition: (lzards.yaml)
---
title: LZARDS NAMS Authorization Demo
variables:
  - name: BASE_URL
    value: http://localhost:3000
    echo: true
  - name: LZARDS_TOKEN
    cmd: ./scripts/getToken.js
    echo: false
setup:
  - comment: Blow away the existing LZARDS env and DB.
    cmd: docker-compose down
  - comment: Run database migrations.
    cmd: docker-compose run api knex migrate:latest
  - comment: Start up the LZARDS service.
    cmd: docker-compose up
    background: true
demo:
  - comment: >
      This is all up and running in SIT, but I'll be demoing locally so that I
      can show what it looks like when you don't have access at first.
  - comment: Hit the status/health check endpoint, which doesn't require authentication.
    cmd: curl -i ${BASE_URL}/api/status
  - comment: Hit the backups listing endpoint, which _does_ require authentication.
    cmd: curl -i ${BASE_URL}/api/backups
  - comment: Try again, with the token we got from the Launchpad API.
    cmd: 'curl -i ${BASE_URL}/api/backups -H "Authorization: Bearer ${LZARDS_TOKEN}"'
  - comment: >
      As we can see in the logs, I was successfully authenticated as the
      svgs-lzards-dev@ndc.nasa.gov service account, but that service account
      _hasn't_ been granted any permissions in LZARDS yet.
  - comment: Adding the service account to the "users" role.
    cmd: ./scripts/addServiceAccount.js svgs-lzards-dev@ndc.nasa.gov user
  - comment: Then re-running the previous request.
    cmd: 'curl -i ${BASE_URL}/api/backups -H "Authorization: Bearer ${LZARDS_TOKEN}"'
  - comment: (There's no backups listed because I wiped out the DB before this demo.)
teardown:
  - cmd: docker-compose stop
*/

/* Theoretical usage:
$ demo setup -f {filename}
# Run setup steps
$ demo run -f {filename}
# Run the demo itself
$ demo teardown -f {filename}
# Run the teardown steps

Maybe filename is optional--default to `godemo.yaml` in the current working
directory?

The setup, demo, and teardown steps are really all the same, and execute the same
way, it's just a way to organize things.

Each demo step should run and then wait for input--hit enter to continue, or
enter a specific keyword/instruction to alter the next step (or exit early).

When setup, run, or teardown are executed, first output the demo commands that
can be entered after each command.
*/

var cli struct {
	Setup    SetupCmd    `kong:"cmd,help='Run the setup steps.'"`
	Run      RunCmd      `kong:"cmd,help='Run the demo.'"`
	Teardown TeardownCmd `kong:"cmd,help='Run the teardown steps.'"`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnMissing())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
