# Install

## MySQL

``` mysql
# Install Mercurial, MYSQL and Go

# Create database "intra_hub":

mysql -uroot
CREATE DATABASE intra_hub;
USE intra_hub;
CREATE TABLE `session` (
	`session_key` char(64) NOT NULL,
	`session_data` blob,
	`session_expiry` int(11) unsigned NOT NULL,
	PRIMARY KEY (`session_key`)
	) ENGINE=MyISAM DEFAULT CHARSET=utf8;
```

## ConfPerso

```
mkdir src/intra-hub/confperso
cd src/intra-hub/confperso
touch conf.go
```

Open conf.go et put that inside :

```
package confperso

// To update for yourself. Don't commit it
const (
    AliasDbName                = "default"
    DatabaseName               = "intra_hub"
    Username                   = "root"
    Password                   = ""

    EmailUsername = "vincent.neel@epitech.eu"
    EmailPassword = "*******"
    EmailHost     = "smtp.epitech.eu"
    EmailHostPort = "587"
)
```

Replace username and password by the value you need.  


## Dependencies

```
go get -v -t ...
go get github.com/beego/bee
```


## Beego path

```
# Add bee to your GOPATH
sudo mv $GOPATH/bin/bee /usr/bin
```

## How to create the database

```
cd src/intra-hub
bee run
# After it has build, quit the app (ctrl + c)
./intra-hub orm syncdb --force=true
```

# How to Run

```
cd src/intra-hub
bee run
```

# Bower Dependencies 

```
# Install npm, bower

cd src/intra-hub/static
bower install

```

## How to load users

Browse to http://localhost:8088/task and run the task named BlowfishCrawler
Or make a post request to http://localhost:8088/task?taskname=blowfishCrawler

# Contribution

- master : branche de production  
- develop : branche principale de dev
- feature-XXX : branche pour ajouter la feature XXX  

Apres avoir fini une feature, la merge (ou mieux, la rebase) sur la branche develop.

# Next features

## Page d'ajout de projet

Nom du projet  
Liste des étudiants (login + promotion)  
Court descriptif  
Nom du manageur assigné  

## Page de modification de projet

R.A.S

## Page de visualisation des projets

Tri par : Promotion(s) impliqué(s), étudiant, état d'avancement, thémes, techno, manageur assigné


# Database, backlog & documentation

https://drive.google.com/drive/u/0/folders/0B_SWVXj-Hf43fmlXX0VhV2QxTVA3U2NCX0RSYmJRblgtXzdkTkJnWWNQY21IT1lQWExiY0E

