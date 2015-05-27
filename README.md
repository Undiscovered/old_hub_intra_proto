# Install

## MySQL

```
Install MYSQL and Go

Create database "intra_hub":

mysql -uroot
> CREATE DATABASE intra_hub;

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
)
```

Replace username and password by the value you need.  


## Dependencies

```
go get github.com/astaxie/beego   
go get github.com/beego/bee
go get github.com/beego/i18n
go get github.com/bitly/go-simplejson
go get golang.org/x/crypto/bcrypt
```

## Beego path

```
# Add bee to your GOPATH
sudo mv $GOPATH/bin/bee /usr/bin
```

# How to Run

```
cd src/intra-hub
bee run
```

# How to create the database

```
cd src/intra-hub
bee run
# After it has build, quit the app (ctrl + c)
./intra-hub orm syncdb --force=true
```

# Gestion des branches

- master : branche de production  
- develop : branche principale de dev
- feature-XXX : branche pour ajouter la feature XXX  

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

