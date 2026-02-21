# GIT

| Action                                   | Commande principale                         |
| ---------------------------------------- | ------------------------------------------- |
| Vérifier l’état                          | `git status`                                |
| Ajout                                    | `git add .`                                 |
| Enregistrer localement                   | `git commit -m "message"`                   |
| Envoyer sur GitHub                       | `git push origin <branche>`                 |
| Mettre à jour depuis GitHub              | `git pull origin <branche>`                 |
| ------------------------ Raccourci de    | `git fetch origin master`                   |
| ---------------------------- et          | `git merge FETCH_HEAD`                      |
| Créer une branche                        | `git checkout -b <branche>`                 |
| Changer de branche                       | `git switch <branche>`                      |
| Créer un tag stable                      | `git tag -a vX.Y -m "description"`          |
| Rétablir un fichier                      | `git restore <fichier>`                     |
| Revenir vers une version antérieure      | `git checkout c99ef87`                      |
| Faire la différence avec la version head | `git diff c99ef87..dev -- shelly-proxy.php` |
| Voir tous les tags                       | `git tag`                                   |

🧩 Exemple complet : tag stable à partir de dev

# 1. Sur la branche de développement

git checkout dev

# 2. Fusion vers master

git checkout master
git pull origin master
git merge dev

# 3. Créer le tag stable

git tag -a v1.3 -m "Version stable issue de dev - amélioration ..."

# 4. Pousser vers GitHub

git push origin master
git push origin v1.3

Le tag v1.3 sera rattaché au même commit que master et visible comme version stable sur GitHub.

# Fil d'ariane

En utilisant la méthode de

```shell
git describe --tags
```

Le numéro de version devient un fil d'Ariane qui relie ton binaire à l'historique de ton code.
Ce que l'utilisateur (ou toi) verra concrètement
Imaginons la chronologie suivante :

- Sur Master : Tu as créé le tag v0.1.0.
- La commande affiche : v0.1.0
- Sur Dev : Tu as ajouté 3 commits depuis ce tag.
- La commande git describe --tags --always affichera : v0.1.0-3-g7a8b9c

Le décodage de cette version "Dev" :

- v0.1.0 : La dernière base stable connue.
- 3 : Le nombre de commits effectués depuis cette base.
- g7a8b9c : Le "g" pour Git + le hash court du commit actuel.

Imagine qu'un utilisateur te dise : "Ton serveur plante avec la version v0.1.0-12-a1b2c3d".

- Tu n'as pas besoin de lui demander quand il a téléchargé le code.
- Tu fais un git checkout a1b2c3d dans ton projet.
- Tu te retrouves exactement dans le même état de code que lui au moment du bug.
