# Install

## Beego

```
go get github.com/astaxie/beego   
go get github.com/beego/bee  
```

## Dep

```
go get golang.org/x/crypto/bcrypt
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

.

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

