# acquisition-backend

![CircleCI](https://circleci.com/gh/TSAP-Laval/acquisition-backend.svg?style=svg)
[![Coverage Status](https://coveralls.io/repos/github/TSAP-Laval/acquisition-backend/badge.svg?branch=master)](https://coveralls.io/github/TSAP-Laval/acquisition-backend?branch=master)
[![CodeFactor](https://www.codefactor.io/repository/github/tsap-laval/acquisition-backend/badge)](https://www.codefactor.io/repository/github/tsap-laval/acquisition-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/TSAP-Laval/acquisition-backend)](https://goreportcard.com/report/github.com/TSAP-Laval/acquisition-backend)
[![Code Climate](https://codeclimate.com/repos/58e82ea5c062d70264000dc5/badges/1beb1535545c1e464028/gpa.svg)](https://codeclimate.com/repos/58e82ea5c062d70264000dc5/feed)
[![Issue Count](https://codeclimate.com/repos/58e82ea5c062d70264000dc5/badges/1beb1535545c1e464028/issue_count.svg)](https://codeclimate.com/repos/58e82ea5c062d70264000dc5/feed)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/81cf4d96c1fc41e6992c22aadca440a5)](https://www.codacy.com/app/laurentlp/acquisition-backend?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=TSAP-Laval/acquisition-backend&amp;utm_campaign=Badge_Grade)

## Procédures pour PostgreSQL
### Sur Mac :
#### Installation :
  ```bash
  $ brew install postgres
  ```
#### Démarrer le serveur :
  ```bash
  $ postgres -D /usr/local/var/postgres
  $ createdb 'whoami'
  ```
  ```bash
  $ cd jusqu'au dossier bin de postgres
  $ psql -h localhost -p 5432 -U postgress NomBD
  ```
 
#### Pour accéder au shell Postgres :
  ```bash
  $ psql <nom de la BD>
  ```
  ```bash
  $ cd jusqu'au dossier bin de postgres
  $ plsql -U Prostgres
  ```
#### Pour importer/exporter un script PG-SQL en invite de commande :
##### Importer
  ```bash
  $ psql <nom de la BD> < fichier.pgsql
  ```
  ```bash
  $ psql -h hostname -p port_number -U username -f your_file.sql databasename  
  bd doit être créé
  ```
##### Exporter
  ```bash
  $ psql <nom de la BD> > fichier.pgsql
  ```
  
  

## Procédures pour Nginx
### Sur Mac :
#### Installation :
  ```bash
  $ brew install nginx
  ```
#### Démarrer le serveur :
  ```bash
  $ sudo nginx
  ```
#### Arrêter le serveur :
  ```bash
  $ sudo nginx -stop
  ```

#### Dossier pour les configurations Nginx :
  ```bash
  $ vi /usr/local/etc/nginx/nginx.conf
  ```

### Procédure pour tester l'API
  ```bash
  $ go test -v -race ./...
  ```
enlever le -v pour un résultat abrégé

### Pour tester le code coverage en local :
  ```bash
  go test -v -cover -race -coverprofile=coverage.out ./api
  ```
  pour voir les résultats et les correctifs à apporter en format html (ouvre le navigateur par défaut)
  ```bash
  $ go tool cover -html=coverage.out
  ```

