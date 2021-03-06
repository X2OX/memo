# [Memo](https://github.com/X2OX/memo)

[![JetBrains Open Source Licenses](https://img.shields.io/badge/-JetBrains%20Open%20Source%20License-000?style=flat-square&logo=JetBrains&logoColor=fff&labelColor=000)](https://www.jetbrains.com/?from=blackdatura)
[![Docker Pulls](https://img.shields.io/docker/pulls/x2ox/memo.svg)](https://hub.docker.com/r/x2ox/memo)
[![Release](https://img.shields.io/github/v/release/x2ox/memo.svg)](https://github.com/X2OX/memo/releases)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

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

