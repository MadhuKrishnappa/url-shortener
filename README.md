
# URL-SHORTENER

Golang-based url-shortener coupled with Docker


## Deployment

To deploy this project run

```bash
  docker-compose build --no-cache && docker-compose up -d
```


## Features

- Basic Golang REST APIs
- Mysql Integration
- Golang & Mysql Images
- Schema creation druing deployment


## Usage

Use below link to export postman collection

[POSTMAN Collection](https://documenter.getpostman.com/view/154281/2s8YzTUNLx)


## Run Locally

Clone the project

```bash
  git clone git@github.com:MadhuKrishnappa/url-shortener.git
```

Go to the project directory

```bash
  cd url-shortener
```

DB Connection changes

```bash
  db/url_mapping_dao.go
```

DB Schema

```bash
  db/migration.sql
```

Start the server

```bash
  go run main.go
```
