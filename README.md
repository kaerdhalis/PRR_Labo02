#PRR - Laboratoire 2 : Exclusion mutuelle

## Lancement et parametrage de l'application

Le parametreage des applications se fait dans le fichier json **config.json** du package config,
il permet de determiner le nombre de processus, leur adresses Ip ainsi que les delais artificiel.
Le champs ArtificialDelay determine le temps en seconde que chaque processus va passer dans la section critique avant de rendre le mutex.
Le lancement de l'application se fait grace au launcher.py qui se trouve a la racine du repo. il suffit de l'executer avec la commande
**python launcher.py**. (à noter que son execution a été testée sur windows mais pas sur mac et si il ne fonctionne pas il suffit de lancer autant de client que précisé dans le fichier avec la commande **go run client.go <id> ** avec les id allant de 0 a n-1)

Une fois les processus lancés,l'utilisateur peut renter differentes commandes dans le terminal pour interagir avec le processus.
la commande **r** permet d'afficher la variable globale, la commande **w \<value>** permet de changer la valeur de la variable et la commande **d \<float>** permet de modifier le delai de transmission des messages du porcessus(par defaut ce delai est mis a 0).
Nous avons pris la decision de limiter le parametrage au fichier json uniquement et de ne pas accepter les parametres en ligne de commande car nous estimons que passer des adresses en parametre est tres fastidieux et peut engendrer rapidement des errreurs.

## Implementation

###
