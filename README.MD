# Example Bot

This is a [GOBBL](https://github.com/calebhiebert/gobbl) facebook messenger bot template, follow the instructions below to get started

## To Run

#### 1. Clone

```
git clone https://github.com/calebhiebert/gobbl-starter-template
```

#### 2. Environment Setup

- Create a new file called `.env`
- Copy the contents of the `env` file into `.env`
- Fill in the `.env` file

#### 3. Dependencies

Install go dependency packages

```bash
go get ./...
```

This starter project comes complete with some database functionality built in. To make the process less cumbersome, a `docker-compose` file is included. Just run `docker-compose up -d` to start up the database.

#### 4. Update Facebook Webhook

If using ngrok for example, url would look like this: https://da8d1354.ngrok.io/webhook

#### 5. That's It!

You should be good to go!

## Project Layout

#### Assets

This template includes static asset serving. Any static assets should be placed in the `/assets/static` folder. The `UStaticAssetURL` function can be used to easily generate full URLs to static assets. For example, if you had an image called `hello.png` placed at `/assets/static/hello.png`, you would call `UStaticAssetURL("hello.png")` to get a full URL.

#### Internationalization

This template includes internationalization support. Language files should be placed in the `/assets/lang` folder. For convenience an `en-US` file is already created for you. Internationalization is based on https://github.com/nicksnyder/go-i18n. Language files should follow the format defined in that repo.

#### Database

This template includes a basic database to track users and user actions within the bot. The database schema is defined in `/assets/schema.sql`, this sql file is executed upon bot startup by default. Code for interacting with the database should go in the `bdb` package. This template uses the https://github.com/upper/db package for database interaction. Each database table should have a corresponding struct placed in it's own file.

#### File Conventions

- Files that contain gobbl handlers should be prefixed with `h-`
- Files that contain gobbl handlers should not include any other functions
- gobbl handler function names should be prefixed with a capital H
- Files that contain helper/utility methods should be prefixed with `u-`
- Utility/helper files should not contain handler functions

#### Environment Variables

- `PAGE_ACCESS_TOKEN` - The facebook page access token for the page the bot is being hosted on
- `SELF_HOST` - The host name of the place where the bot is hosted (probably http://localhost:8080 for development)
- `ENABLE_DEV_TOOLS` - Enables some utility routes within the bot, see h-dev.go
- `LOG_LEVEL` - The gobbl logging level, should be a numerical value
- `LOG_COLORS` - Pretty colors
