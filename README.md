# RSS-Feed-Bot

## Usage
The original bot is placed as [@TheDispatchBot](https://t.me/TheDispatchBot) in Telegram (but most likely, it is disabled due to the lack of a server :( ). So, it is preferable to run it yourself.

### Prerequisites
- Database

### Command Line Arguments
| Flag      | Description, example                                                                                                                                                                                                                                                                           |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-token`  | Token of the telegram bot that is given by [@BotFather](https://t.me/BotFather)                                                                                                                                                                                                                |
| `-db-url` | Database URL for connection. To run from Docker with database on host machine use `host.docker.internal` as host address. For example, `postgres://username:password@host.docker.internal:port/database`.                                                                                      |
| `-upd`    | Feed update interval in duration format (a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \"300ms\", \"-1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".) Default values: `12 minutes`. |
