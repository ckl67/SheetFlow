# Création branche

git branch dev

# Aller sur la branche

git checkout dev

# Fusionner dev dans master avant de taguer

git checkout master
git merge dev

| Action                                   | Commande principale                         |
| ---------------------------------------- | ------------------------------------------- |
| Vérifier l’état                          | `git status`                                |
| Enregistrer localement                   | `git add .` + `git commit -m "message"`     |
| Envoyer sur GitHub                       | `git push origin <branche>`                 |
| Mettre à jour depuis GitHub              | `git pull origin <branche>`                 |
| Créer un tag stable                      | `git tag -a vX.Y -m "description"`          |
| Rétablir un fichier                      | `git restore <fichier>`                     |
| Revenir vers une version antérieure      | `git checkout c99ef87`                      |
| Faire la différence avec la version head | `git diff c99ef87..dev -- shelly-proxy.php` |
| Créer une branche                        | `git checkout -b <branche>`                 |
| Changer de branche                       | `git switch <branche>`                      |
| Voir tous les tags                       | `git tag`                                   |
