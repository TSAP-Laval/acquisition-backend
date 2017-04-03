# acquisition-backend

![CircleCI](https://circleci.com/gh/TSAP-Laval/acquisition-backend.svg?style=svg)
[![Coverage Status](https://coveralls.io/repos/github/TSAP-Laval/acquisition-backend/badge.svg?branch=master)](https://coveralls.io/github/TSAP-Laval/acquisition-backend?branch=master)

## Procédures pour PostgreSQL
### Sur Mac :
#### Installation :
  ```
  $ brew install postgres
  ```
#### Démarrer le serveur :
  ```
  $ postgres -D /usr/local/var/postgres
  $ createdb 'whoami'
  ```
   ```
  $ cd jusqu'au dossier bin de postgres
  $ psql -h localhost -p 5432 -U postgress NomBD
  ```
 
#### Pour accéder au shell Postgres :
  ```
  $ psql <nom de la BD>
  ```
   ```
  $ cd jusqu'au dossier bin de postgres
  $ plsql -U Prostgres
  ```
#### Pour importer/exporter un script PG-SQL en invite de commande :
##### Importer
  ```
  $ psql <nom de la BD> < fichier.pgsql
  ```
  ```
  $ psql -h hostname -p port_number -U username -f your_file.sql databasename  
  bd doit être créé
  ```
##### Exporter
  ```
  $ psql <nom de la BD> > fichier.pgsql
  ```
  
  

## Procédures pour Nginx
### Sur Mac :
#### Installation :
  ```
  $ brew install nginx
  ```
#### Démarrer le serveur :
  ```
  $ sudo nginx
  ```
#### Arrêter le serveur :
  ```
  $ sudo nginx -s stop
  ```
#### Dossier pour les configurations Nginx :
  `vi /usr/local/etc/nginx/nginx.conf`
  
## Procédure pour tester l'API
  `$ go test -v -race ./...`
  enlever le -v pour un résultat abrégé

# Pour tester le code coverage en local :
  `go test -v -cover -race -coverprofile=coverage.out ./api`
  pour voir les résultats et les correctifs à apporter
  `$ go tool cover -html=coverage.out`

