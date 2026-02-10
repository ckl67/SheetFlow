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
