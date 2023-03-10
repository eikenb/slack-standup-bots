# Lambda + DynamoDB based standup bot for Slack

See Makefile for details on building/testing/etc. Though, assuming you have
access, all you need to get it going is..

    make all

Then take the URL it outputs at the end and use in Slack config.

### Slack Setup

While Logged into Slack with Admin privileges.

1. Goto https://api.slack.com/apps
2. Click "Create New App"
3. Click "Slash Commands" (under "Add features and functionality")
4. Fill out new command form:
  - Command: /up
  - Request URL: use URL output from `make all` command above
  - Short Description: standup bot
  - Usage Hint: /up help
5. Click "Save"
6. Click "Basic Information"
7. Scroll down to "App Credentials"
8. Write down "Verification Token"
  - Encrypt with KMS (see below) and put in template.yaml
  - run `make all` again to update lambda (URL won't change)

### KMS

To encrypt your token, run this command.

    aws kms encrypt --key-id alias/upbot --plaintext SECRET_TOKEN --output text

Note 'alias/upbot' is the current key in use in our setup. If you are setting this up fresh, switching accounts, etc. then see below.

If you don't have an key w/ an ALIAS yet you should first create a new key,
then create the alias to it and use that as ALIAS.

    aws kms create-key --description "key for slack bot"
    aws kms --alias-name ALIAS --target-key-id [KeyId from above]
