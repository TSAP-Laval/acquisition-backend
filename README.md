# acquisition-backend

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
#### Pour importer/exporter un script PGSQL en invite de commande :
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
#### Dossier pour le fichier de configurations Nginx :
  `vi /usr/local/etc/nginx/nginx.conf`
  
