Cortex-Analyzers in Go
----------------------

This repo contains sources for [Cortex-Analyzers](https://github.com/TheHive-Project/Cortex-Analyzers) written in Go

## List of analyzers

* [Dor](https://github.com/ilyaglow/dor) analyzer - domank ranking
* [HaveIBeenPwned](https://haveibeenpwned.com) analyzer - check an email for breaches and pastes
* [HackedEmails](https://hacked-emails.com) analyzer - check an email for breaches and pastes, sometimes better, sometimes worse than hibp
* [BadPackets](https://mirai.badpackets.net) analzyer - check if IP address has been seen in Mirai-like Botnet list by badpackets.net

## How to build

If you have Go installed on your host you can use a simple make command:

```
make
```

You can also use docker to make builds:
```
make docker-build
```
