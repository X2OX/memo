# memo 随手笺

## Configuration

```json
{
    "database":"db.db",
    "domain":"https://domain.com",
    "listen_addr":":8088",
    "data_folder":"/data/memo/",
    "log_level":"error",
    "telegram_id":123456789,
    "telegram_token":"123456789:abc",
    "telegram_webhook":"/telegram/webhook",
    "token":{
        "auto_update":0,
        "preview":10,
        "view":0,
        "share":1
    }
}
```

- `database` data source name, support `SQLite3` and `PostgreSQL`
- `domain` used by the webhook, preview, share
- `listen_addr` program listening address
- `data_folder` work folder
- `log_level` log level
- `telegram_id` your telegram id, isn't username
- `telegram_token` Bot's token
- `telegram_webhook` webhook path, switch randomly will cause the message to be lost
- `token.auto_update` how many minutes to update the token, Disable when zero
- `token.preview` the effective minutes of the preview link
- `token.view` the effective minutes of the view link
- `token.share` the effective minutes of the share link

