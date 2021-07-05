# [随手笺](https://github.com/X2OX/memo)

[![JetBrains Open Source Licenses](https://img.shields.io/badge/-JetBrains%20Open%20Source%20License-000?style=flat-square&logo=JetBrains&logoColor=fff&labelColor=000)](https://www.jetbrains.com/?from=blackdatura)
[![Docker Pulls](https://img.shields.io/docker/pulls/x2ox/memo.svg)](https://hub.docker.com/r/x2ox/memo)
[![Release](https://img.shields.io/github/v/release/x2ox/memo.svg)](https://github.com/X2OX/memo/releases)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## 配置

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

- `database` DSN，目前只支持 `SQLite3` 和 `PostgreSQL`
- `domain` 机器人 Webhook 及访问会用到，和监听地址不同
- `listen_addr` 程序监听的地址及端口
- `data_folder` 工作目录
- `log_level` Log 记录的级别
- `telegram_id` 你的 Telegram ID，不是用户名
- `telegram_token` Bot 的 token
- `telegram_webhook` Webhook path 不需要加域名，频繁切换模式可能会丢失消息
- `token.auto_update` 密钥自动更新时间「分钟」
- `token.preview` 预览链接的有效期「分钟」
- `token.view` 阅读链接的有效期「分钟」
- `token.share` 分享链接的有效期「分钟」

