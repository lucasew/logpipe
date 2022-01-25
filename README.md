# logpipe
A simple python script to be used as a log pipe. From your services to a external sink like Discord or Telegram.

## How to use
- Setup a configuration file. Each section defines settings for a module.
- There are two types of actors: sources and sinks: Sources emit text lines, sinks deal with the lines.
- Reference the configuration file using the `-c` flag

## Already working
- Sources
  - Basic systemd journalctl interactions: filters might be added later
- Sinks
  - Terminal: just output to stdout
  - Telegram: send to `chat_id` using `token`
    - `TELEGRAM_TOKEN` and `TELEGRAM_CHAT` environment variables are supported too as fallback!
  - Discord: send to a given chat using the generated webhook URL
    - `DISCORD_WEBHOOK` environment variable is supported as fallback!
  
## Limitations
- Async support. So far everything is synchronous.
- Rate limit: bursts of lines can trigger ratelimit limitations of the services.
- Each line turns into a request to a external service: no batching mechanism yet.

## Dependencies
- This program uses only Python 3.x and it's standard library and was tested on Python 3.9.9.
