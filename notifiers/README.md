# Notifications

## Slack

First create a Slack App with `chat:write` and `chat:write.public` (optional) scopes

Then specify 1 or more slack apps (with `name` and `auth-token` bot user OAuth token) and channels (`name` and `channel-id` for each)

**Configuration:**
- `auth-token`
    - type: `string`
    - required: **yes**
- `channel-id`
    - type: `string`
    - required: **yes**

**Sample configuration:**
```yaml
notifications:
  
    ... [reducted]

    slack:
        apps:                   # list of slack apps
            - name: "go-alive-test-bot"
              token: "add bot user oauth token here"
        channels:               # list of slack channels
            - name: "go-alive-test-group"
              channelid: 'CHANNEL ID'
```  
