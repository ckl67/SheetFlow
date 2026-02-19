package forms

import (
	"errors"
	"mime/multipart"
	"strings"
)

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
// ValidateForm() est une méthode de la struct UploadRequest qui est actuellement un no-op (ne fait rien).
// Cependant, elle est définie pour permettre l'ajout de toute validation personnalisée du formulaire à l'avenir sans avoir à modifier le code qui appelle cette méthode.
// Par exemple, si tu veux ajouter une validation pour vérifier que le champ SheetName n'est pas vide, tu pourrais implémenter ValidateForm() comme ceci :
//
//	func (req *UploadRequest) ValidateForm() error {
//	    if strings.TrimSpace(req.SheetName) == "" {
//	        return errors.New("SheetName is required")
//	    }
//	    return nil
//	}
func (req *UploadRequest) ValidateForm() error {
	if req.File == nil {
		return errors.New("file is required")
	}

	if strings.TrimSpace(req.Composer) == "" {
		return errors.New("composer is required")
	}

	if strings.TrimSpace(req.SheetName) == "" {
		return errors.New("sheet name is required")
	}

	if req.File.Size > 10<<20 { // 10MB
		return errors.New("file too large")
	}

	if !strings.HasSuffix(strings.ToLower(req.File.Filename), ".pdf") {
		return errors.New("only PDF files are allowed")
	}

	return nil
}
