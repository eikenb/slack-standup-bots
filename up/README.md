Up
--

Simple slack stand-up bot.

    Usage: @up command (or in private channel)
    Commands:
        standup <what you did yesterday, doing today, etc.>
            How you enter your standup.
            Echo's it back in main channel if you do it in private channel.
            Free form. All 1 line.
        status
            Outputs everyone's most recent standup entries.
    Eg.
    @up yesterday I worked on bug #1; today I worked on bug #2; no blockers

Deployment
----------

Requires 2 environment variables.

- SLACK_TOKEN - your slack bot api token.
- REDIS_HOST  - hostname:port of redis server.

Logs to STDOUT.
