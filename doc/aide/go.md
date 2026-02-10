# Rôle de go.mod

go.mod dit à Go : « Voici le module, voici ma version de Go, et voici mes dépendances
Même localement, Go regarde le module pour résoudre les imports et vérifier les versions.

# Rôle de go mod tidy

go mod tidy

fait 2 choses principales :
Ajoute les dépendances manquantes :
Si un package externe (ex : github.com/gin-gonic/gin) n’est pas dans go.mod, tidy l’ajoute automatiquement.
Supprime les dépendances inutilisées :
Si go.mod contient des packages qui ne sont plus importés dans le code, tidy les supprime.
C’est un nettoyage et synchronisation du fichier go.mod par rapport au code réel.

# config linux

Pensez à mettre cette ligne dans un fichier de configuration du shell ~/.bashrc

export PATH="$PATH:/usr/local/go/bin"

Puis recharge : source ~/.bashrc
