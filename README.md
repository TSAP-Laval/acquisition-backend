# acquisition-backend

[![CircleCI](https://circleci.com/gh/TSAP-Laval/acquisition-backend.svg?&maxAge=3600&style=shield)]
[![Coverage Status](https://coveralls.io/repos/github/TSAP-Laval/acquisition-backend/badge.svg?maxAge=3600&branch=master)](https://coveralls.io/github/TSAP-Laval/acquisition-backend?branch=master)
[![CodeFactor](https://www.codefactor.io/repository/github/tsap-laval/acquisition-backend/badge?maxAge=3600)](https://www.codefactor.io/repository/github/tsap-laval/acquisition-backend)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/81cf4d96c1fc41e6992c22aadca440a5?maxAge=3600)](https://www.codacy.com/app/laurentlp/acquisition-backend?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=TSAP-Laval/acquisition-backend&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/TSAP-Laval/acquisition-backend?maxAge=3600)](https://goreportcard.com/report/github.com/TSAP-Laval/acquisition-backend)
[![Code Climate](https://codeclimate.com/github/TSAP-Laval/acquisition-backend/badges/gpa.svg?maxAge=3600)](https://codeclimate.com/github/TSAP-Laval/acquisition-backend)
[![Issue Count](https://codeclimate.com/github/TSAP-Laval/acquisition-backend/badges/issue_count.svg?maxAge=3600)](https://codeclimate.com/github/TSAP-Laval/acquisition-backend)

## Procédures pour PostgreSQL

### Sur Mac

#### Installation

  ```bash
  brew install postgres
  ```

#### Démarrer le serveur

  ```bash
  postgres -D /usr/local/var/postgres
  createdb 'whoami'
  ```

  ```bash
  cd jusqu'au dossier bin de postgres
  psql -h localhost -p 5432 -U postgress NomBD
  ```

#### Pour accéder au shell Postgres

  ```bash
  psql <nom de la BD>
  ```

  ```bash
  cd jusqu'au dossier bin de postgres
  plsql -U Prostgres
  ```

#### Pour importer/exporter un script PG-SQL en invite de commande

##### Importer

  ```bash
  psql <nom de la BD> < fichier.pgsql
  ```
  ```bash
  psql -h hostname -p port_number -U username -f your_file.sql databasename
  bd doit être créé
  ```

##### Exporter

  ```bash
  psql <nom de la BD> > fichier.pgsql
  ```

## Procédures pour Nginx

  ```bash
  brew install nginx
  ```

### Démarrer le serveur nginx

  ```bash
  sudo nginx
  ```

### Arrêter le serveur

  ```bash
  sudo nginx -stop
  ```

### Dossier pour les configurations Nginx

  ```$ vi /usr/local/etc/nginx/nginx.conf```

### Procédure pour tester l'API

  ```bash
  go test -v ./api
  ```
enlever le -v pour un résultat abrégé (plus rapide)

### Pour tester le code coverage en local

  ```bash
  go test -v -cover -coverprofile=coverage.out ./api
  ```
  pour voir les résultats et les correctifs à apporter en format html (ouvre le navigateur par défaut)

  ```bash
  go tool cover -html=coverage.out
  ```
