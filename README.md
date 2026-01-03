# Projet Final : API pour une Clinique Vétérinaire


## Table des matières

- [Lancement du projet](#lancement-du-projet)
- [Les Routes](#les-routes)
  - [Chat](#chat)
  - [Visite](#visite)
  - [Traitement](#traitement)
  - [Utilisateur](#utilisateur)
  - [Authentification](#authentification)
- [Architecture](#architecture)


## Lancement du projet

Aller à la racine et executer la commande :

```
go run .
```

L'API serait alors disponible sur **http://localhost:8081/api/v1/vet**

Une documentation Swagger complete est aussi disponible sur **http://localhost:8081/swagger/index.html**

## Les Routes

### Chat
<details>
<summary><strong>Voir les routes chat</strong></summary>

| Méthode | Endpoint | Description | Auth |
|---------|---------|------------|------|
| POST    | /cats | Ajouter un chat | admin |
| GET     | /cats | Récupérer tous les chats | all |
| GET     | /cats/{id} | Récupérer un chat par son ID | all |
| GET     | /cats/{id}/history | Historique des visites du chat | all |
| PUT     | /cats/{id} | Modifier un chat | admin |
| DELETE  | /cats/{id} | Supprimer un chat | admin |

</details>

### Visite
<details>
<summary><strong>Voir les routes visite</strong></summary>

| Méthode | Endpoint | Description | Auth |
|---------|---------|------------|------|
| POST    | /visits | Ajouter une visite | admin |
| GET     | /visits | Récupérer toutes les visites | all |
| GET     | /visits/{id} | Récupérer une visite par son ID | all |
| PUT     | /visits/{id} | Modifier une visite | admin |
| DELETE  | /visits/{id} | Supprimer une visite | admin |

</details>

### Traitement
<details>
<summary><strong>Voir les routes traitement</strong></summary>

| Méthode | Endpoint | Description | Auth |
|---------|---------|------------|------|
| POST    | /treatments | Ajouter un traitement | admin |
| GET     | /treatments | Récupérer tous les traitements | all |
| GET     | /treatments/{id} | Récupérer un traitement par son ID | all |
| GET     | /treatments/{id}/history | Récupérer les traitements associés à une visite | all |
| PUT     | /treatments/{id} | Modifier un traitement | admin |
| DELETE  | /treatments/{id} | Supprimer un traitement | admin |

</details>

### Utilisateur
<details>
<summary><strong>Voir les routes utilisateur</strong></summary>

| Méthode | Endpoint | Description | Auth |
|---------|---------|------------|------|
| POST    | /users | Ajouter un utilisateur | admin |
| GET     | /users | Récupérer tous les utilisateurs | all |
| GET     | /users/{id} | Récupérer un utilisateur par son ID | all |
| PUT     | /users/{id} | Modifier un utilisateur | admin |
| DELETE  | /users/{id} | Supprimer un utilisateur | admin |

</details>

### Authentification
<details>
<summary><strong>Voir les routes authentification</strong></summary>

#### Recevoir un access et refresh token pour un utilisateur
- **POST** `/users/login` (admin only)
On envoie dans le body email plus mot de passe et on reçoit les deux tokens.

Requête
```
{
  "user_email": "user@example.com",
  "user_password": "password123"
}
```

Réponse
```
{
  "access_token": "...",
  "refresh_token": "..."
}
```

Les tokens sont indispansable pour faire des requetes sur toutes les routes de l'API à l'exception de routes "Users". C'est une sécurité supplémentaire.

#### Recevoir un nouvelle access token
- **POST** `/users/refresh` (admin only)
On envoie dans le body le refresh token et on reçoit un nouvelle access token.

Le refresh token permet de redemander à l'API un nouvelle access token sans avoir à refaire une connexion.

</details>

## Architecture

<details>
<summary><strong>Voir l'arborescence</strong></summary>

```
    ┬ VET-CLINIC-API
    ├───┬ config
    │   └──── config.go
    ├───┬ database
    │   ├──── dbmodel
    │   │       ├──── cat.go
    │   │       ├──── treatment.go
    │   │       ├──── user.go
    │   │       └──── visit.go
    │   └──── database.go
    │
    ├───┬ docs
    │   ├──── docs.go
    │   ├──── swagger.json
    │   └──── swagger.yaml
    │
    ├───┬ pkg
    │   ├───── authentication
    │   │       ├──── jwt.go
    │   │       └──── middleware.go
    │   │
    │   ├───── cat
    │   │       ├──── controller.go
    │   │       └──── routes.go
    │   ├───── models
    │   │       ├──── cat.go
    │   │       ├──── token.go
    │   │       ├──── treatment.go
    │   │       ├──── user.go
    │   │       └──── visit.go
    │   ├───── treatment
    │   │       ├──── controller.go
    │   │       └──── routes.go
    │   ├───── user
    │   │       ├──── controller.go
    │   │       └──── routes.go
    │   └───── visit
    │           ├──── controller.go
    │           └──── routes.go
    │
    ├──── .env
    ├──── .env.example
    ├──── .gitignore
    ├──── go.mod
    ├──── go.sum
    ├──── main.go
    └──── README.md
```
</details>