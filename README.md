Cortex-Analyzers in Go
----------------------

This repo contains sources for [Cortex-Analyzers](https://github.com/TheHive-Project/Cortex-Analyzers) written in Go

## List of analyzers

* [Dor](https://github.com/ilyaglow/dor) analyzer - domank ranking
* HaveIBeenPwned analyzer - check an email for breaches and pastes

## How to build

If you have Go installed on your host you can use a simple make command:

```
make
```

If you don't, no worries:
```
make docker-build
```
