#PRR - Laboratoire 2 : Exclusion mutuelle

## Lancement et parametrage de l'application

Le parametreage des applications se fait dans le fichier json **config.json** du package config,
il permet de determiner le nombre de processus, leur adresses Ip ainsi que les delais artificiel.
Le champs ArtificialDelay determine le temps en seconde que chaque processus va passer dans la section critique avant de rendre le mutex.
Le lancement de l'application se fait grace au launcher.py qui se trouve a la racine du repo. il suffit de l'executer avec la commande
**python launcher.py**. (à noter que son execution a été testée sur windows mais pas sur mac et si il ne fonctionne pas il suffit de lancer autant de client que précisé dans le fichier avec la commande **go run client.go id** avec les id allant de 0 a n-1)

Une fois les processus lancés,l'utilisateur peut renter differentes commandes dans le terminal pour interagir avec le processus.
la commande **r** permet d'afficher la variable globale, la commande **w \<value>** permet de changer la valeur de la variable et la commande **d \<float>** permet de modifier le delai de transmission des messages du porcessus(par defaut ce delai est mis a 0).
Nous avons pris la decision de limiter le parametrage au fichier json uniquement et de ne pas accepter les parametres en ligne de commande car nous estimons que passer des adresses en parametre est tres fastidieux et peut engendrer rapidement des errreurs.

## Implementation

###Client
Module s'occupant de la partie client du processus,c'est par la que le process se lance.Il lance le mutex et le network, attend que celui cii soit operationnel puis permet au client de rentrer les commandes grace a une boucle infinie, il s'occupe egalement de l'affichage de la valeur partagée.
Il prend l'id du processus au paramètre et le passe au mutex.

###Mutex
Implemente l'algorithme de Carvalho et Roucairol pour gerer l'exclusion mutuel ainsi que l'envoie des messages aux differents processus.
La boucle principale est basée sur un select qui attend les signaux du client ou les messages des autres process et permet ainsi une exclusion mutelle.
La boucle principale demarre apres que le processus ait verfié que tous les autres processus ont finit leur setup grace a un ping sur leurs adresses.

###Connection
Package gerant la communication TCP(envoie et reception des messages) et transmet les messages reçus au mutex.

###Config
Package gerant la configuration generale de l'application(nombre de process, leurs adresses et le temps en section crittique) en lisant le fichier config.json

## Ce qui reste a faire
L'application fonctionne et remplie les conditions du laboratoire,cependant des moifications peuvent y etre apporté.
En effet nous pourrions rendre le package network plus generique pour le reutiliser plus simplement dans un autre projet, par exemple en faisant un package tcp n'acceptant que des tableaux de bytes en parametres et un autre package convertissant les messages en binaire  avant de le passer au module gerant le TCP.
Nous aurions egalement pu gerer le cas ou un processus est stoppé(meme si le resueau est supposé fiable) pour que les autres process enleve celui-ci de la liste des process.
Le manque de tests unitaires est aussi une issue.
Un autre point que nous aurions pu ameliorer est le lancement de l'application en faisant des scripts séparés pour les différents OS.
