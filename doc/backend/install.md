# Rappel

GO : langage + compilateur
GoLand : logiciel pour écrire du code

# Supprimer ancienne version Go et Goland

sudo apt remove golang-go golang
sudo rm -rf /usr/local/go

# installation go

_Il est déconseillé d'utiliser sudo snap install go_

wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

vérification
go version
