# Projet Final : API pour une Clinique Vétérinaire


## Lancement du projet

Aller à la racine et executer la commande :

```
go run .
```

L'API serait alors disponible sur **http://localhost:8081/api/v1/vet**

Une documentation Swagger est aussi disponible sur **http://localhost:8081/swagger/index.html**

## Les Routes

#### <span>Chat</span>

Ajoute un chat
- <span><strong style="color:yellow">Post</strong></span> /cats

Récupère toutes les chats
-  <span><strong style="color:green">Get</strong></span> /cats

Récupère un chat en fonction de son Id
- <span><strong style="color:green">Get</strong></span> /cats/{id}

Récupère un chat et son historique de viste qui elle même contienne les traitements
- <span><strong style="color:green">Get</strong></span> /cats/{id}/history

Modifie un chat
- <span><strong style="color:blue">Put</strong></span> /cats/{id}

Supprime un chat
- <span><strong style="color:red">Delete</strong></span> /cats/{id}


#### <span>Visit</span>

Ajoute une visit
- <span><strong style="color:yellow">Post</strong></span> /visits

Récupère toutes les visits
-  <span><strong style="color:green">Get</strong></span> /visits

Récupère une visits en fonction de son Id
- <span><strong style="color:green">Get</strong></span> /visits/{id}

Modifie une visit
- <span><strong style="color:blue">Put</strong></span> /visits/{id}

Supprime une visit
- <span><strong style="color:red">Delete</strong></span> /visits/{id}


#### <span>Traitement</span>

Ajoute un Traitement
- <span><strong style="color:yellow">Post</strong></span> /treatments

Récupère toutes les Traitements
-  <span><strong style="color:green">Get</strong></span> /treatments

Récupère un Traitement en fonction de son Id
- <span><strong style="color:green">Get</strong></span> /treatments/{id}

Récupère les Traitements associé à une visit
- <span><strong style="color:green">Get</strong></span> /treatments/{id}/history

Modifie un Traitement
- <span><strong style="color:blue">Put</strong></span> /treatments/{id}

Supprime un Traitement
- <span><strong style="color:red">Delete</strong></span> /treatments/{id}

## Architecture

```
    ┬ VET-CLINIC-API
    ├───┬ config
    │   └──── config.go
    ├───┬ database
    │   ├──── dbmodel
    │   │       ├──── cat.go
    │   │       ├──── treatment.go
    │   │       └──── visit.go
    │   └──── database.go
    │
    ├───┬ docs
    │   ├──── docs.go
    │   ├──── swagger.json
    │   └──── swagger.yaml
    │
    ├───┬ pkg
    │   ├───── cat
    │   │       ├──── controller.go
    │   │       └──── routes.go
    │   ├───── models
    │   │       ├──── cat.go
    │   │       ├──── treatment.go
    │   │       └──── visit.go
    │   ├───── treatment
    │   │       ├──── controller.go
    │   │       └──── routes.go
    │   └───── visit
    │           ├──── controller.go
    │           └──── routes.go
    │
    ├──── go.mod
    ├──── go.sum
    ├──── main.go
    └──── README.md
```


## Amélioration

Ajout de l'authentification avec JWT
