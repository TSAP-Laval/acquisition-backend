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
#### Pour accéder au shell Postgres :
  ```
  $ psql <nom de la BD>
  ```
#### Pour importer/exporter un script PGSQL en invite de commande :
##### Importer
  ```
  $ psql <nom de la BD> < fichier.pgsql
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
  
