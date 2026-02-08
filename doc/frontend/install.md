# Node

Supprimer complètement Node.js
sudo apt remove --purge nodejs
sudo apt autoremove

## Installer nvm

Permet plus facilemnet de gérer les versions nodejs

curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
source ~/.bashrc

nvm install 20
nvm use 20

# Création du projet avec Vite + React

npm create vite@latest SheetFlow

npm create vite@latest SheetFlow
◇ Package name:
│ sheetflow
◇ Select a framework:
│ React
◇ Select a variant:
│ JavaScript
◇ Use rolldown-vite (Experimental)?:
│ No
◇ Install with npm and start now?
│ Yes

npm install react-router-dom

# Visual Code

## Extensions :

ES7+ React/Redux/React-Native snippets
ESLint Détecte : hooks mal utilisés
Prettier - Code formatter

## Configuration

Fichier .vscode/settings.json

# Git

cd SheetFlow
git init
git config --global user.email "Vous@exemple.com"
git config --global user.name "Votre Nom"
git add .

git commit -m "Initial commit"
git remote add origin https://github.com/ckl67/SheetFlow.git
git push -u origin master
