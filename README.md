## Сервис сокращения ссылок // URL shortener
[![Build status](https://github.com/imotkin/shortener/actions/workflows/build.yml/badge.svg)](https://github.com/github.com/imotkin/shortener/actions/workflows/build.yml)
[![Report card](https://goreportcard.com/badge/github.com/imotkin/shortener)](https://goreportcard.com/report/github.com/imotkin/shortener)

### Запуск приложения
#### 1. С помощью Go

В корневой директории проекта выполнить команду
```shell
go run cmd/main.go
```
#### 2. C помощью Docker
* Создать образ для приложения с помощью команды `docker build`
```shell 
docker build -t shortener .
```
* Запустить контейнер Docker c помощью команды `docker run`
```shell
docker run -dp 5000:5000 shortener
```

После запуска сервис будет доступен по адресу `localhost:5000`

> **Указание:**
> для конфигурации приложения используется файл `config.toml`, в котором указываются хост и порт для сервера приложения,
> а также наименование для базы данных SQLite.

Используемые технологии:
* [Go](https://go.dev) (1.22), [шаблоны Go](https://pkg.go.dev/html/template) для HTML
* базовые стили CSS и скрипты JavaScript
* [драйвер базы данных](https://modernc.org/sqlite) для работы с SQLite
* [goose](https://github.com/pressly/goose) для миграций базы данных
* [библиотека](https://github.com/BurntSushi/toml) для работы с TOML для файла конфигурации

<details>
    <summary>Скриншоты веб-приложения (клик)</summary>
    <br/>
    <img src="https://github.com/imotkin/shortener/blob/main/images/1.png"  alt="Screenshot of main page"/>
    <br/>
    <br/>
    <img src="https://github.com/imotkin/shortener/blob/main/images/2.png"  alt="Screenshot of main page"/>
    <br/>
    <br/>
    <img src="https://github.com/imotkin/shortener/blob/main/images/3.png"  alt="Screenshot of main page"/>
</details>

