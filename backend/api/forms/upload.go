package forms

import "mime/multipart"

// Structure qui sera uilisée par dans la fonction UploadFile() à travers c.ShouldBind(&uploadForm);
// ShouldBind :
//    Parse la requête selon le Content-Type
//    Remplit la struct
//    Retourne une erreur uniquement si le binding échoue
// 		⚠️ Mais il ne rend pas les champs obligatoires par défaut.
//
// Un champ est obligatoire seulement si tu ajoutes :binding:"required"
// Exemple : SheetName string `form:"sheetName" binding:"required"`

// Composer string `form:"composer"`
// Signifie :
//     Gin va chercher un champ nommé "composer" dans le multipart/form-data
//		 Si il trouve un champ "composer", il va convertir sa valeur en string
//		 Si la conversion réussit, il va assigner la valeur convertie à uploadForm.Composer

// Pourquoi *multipart.FileHeader ?
// Parce que c'est la structure que Gin utilise pour représenter les fichiers uploadés
// dans un formulaire multipart/form-data.
// Elle contient des informations sur le fichier, comme son nom, sa taille, son type MIME,
// et une méthode pour ouvrir le fichier et lire son contenu.
// En utilisant *multipart.FileHeader, tu peux facilement accéder au fichier uploadé et le traiter
// dans ta fonction UploadFile().
// En résumé, la struct UploadRequest est conçue pour recevoir les données d'un formulaire d'upload de fichier,
// avec des champs pour le fichier lui-même, le nom du compositeur, le nom de la partition,
// la date de sortie, les catégories, les tags et un texte d'information.
// Le champ File est un pointeur vers multipart.FileHeader,
// ce qui permet de gérer facilement les fichiers uploadés dans Gin.

// File *multipart.FileHeader `form:"uploadFile"`
// → implique obligatoirement :
// multipart/form-data

type UploadRequest struct {
	File            *multipart.FileHeader `form:"uploadFile"`
	Composer        string                `form:"composer"`
	SheetName       string                `form:"sheetName"`
	ReleaseDate     string                `form:"releaseDate"`
	Categories      string                `form:"categories"`
	Tags            string                `form:"tags"`
	InformationText string                `form:"informationText"`
}

// Currently a no-op but enables us to add any custom form validation in without having to change any calling code.

func (req *UploadRequest) ValidateForm() error {
	return nil
}
