# Installation


# Fonctionnalités

## Page d'ajout de projet

Nom du projet  
Liste des étudiants (login + promotion)  
Court descriptif  
Nom du manageur assigné  

## Page de modification de projet

R.A.S

## Page de visualisation des projets

Tri par : Promotion(s) impliqué(s), étudiant, état d'avancement, thémes, techno, manageur assigné


# BDD

A venir : Modélisation de la BDD. A faire en amont si la techno est séléctionnée.

# Backup

Mise en place du script
`cp backup/backup.sh > /bin`
Mise en place du crontab
`crontab -l | cat backup/crontab.txt  | crontab -`