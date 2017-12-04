# Up

Simple slack stand-up bot.

## Usage
@up command (or in private channel)

### Commands:

#### stand <what you did yesterday, doing today, etc.>
How you enter your standup.  
Echo's it back in main channel if you do it in private channel.  
Free form. All 1 line.  

#### show
Outputs everyone's most recent standup entries.

#### help
Help output.

### Eg.
@up stand yesterday I worked on bug #1; today I worked on bug #2; no blockers

## Deployment

Requires 2 environment variables.

- SLACK_TOKEN - your slack bot api token.
- REDIS_HOST  - hostname:port of redis server

Logs to STDOUT.

If you are deploying using Docker and using the Dockerfile. Then you need
create a 'runuprc' file with the 2 environment variables exported in it. Eg.

    export SLACK_TOKEN=thisisnotarealtoken
    export REDIS_HOST=localhost:6379

Then update the Changelog (the version is derived from this file).

Then...
    make image
    make push
    make deploy

If first time.
To make the repo...
    make repo
To create-stack instead of update-stack...
    make init

