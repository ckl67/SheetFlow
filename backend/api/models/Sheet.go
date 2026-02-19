package models

import (
	"backend/api/config"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Sheet modéle de données pour les partitions musicales
// distinguer deux couches orthogonales :
// gorm:"..." → mapping base de données (ORM → schéma SQL)
// json:"..." → mapping API HTTP (struct Go ↔ JSON)

// REMARQUE
// 	Tags    pq.StringArray `gorm:"type:text[]" json:"tags"`
// 		pq.StringArray  est uniquement utilisé dans postgreSQL, il est mappé à un champ de type text[] dans la base de données,
// 		et il est sérialisé/désérialisé en JSON comme un tableau de chaînes de caractères.
// Afin de garder la compatibilité avec d'autres bases de données, SQLITE, MySQL on va stocker les tags et catégories sous forme de chaînes JSON
// dans la base de données
// Tags et Categories contiennent maintenant :
// ["Classical", "Piano"]

type Sheet struct {
	SafeSheetName   string `gorm:"primary_key" json:"safe_sheet_name"`
	SheetName       string `json:"sheet_name"`
	SafeComposer    string `json:"safe_composer"`
	Composer        string `json:"composer"`
	ReleaseDate     time.Time
	PdfUrl          string    `json:"pdf_url"`
	UploaderID      uint32    `gorm:"not null" json:"uploader_id"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Tags            string    `gorm:"type:TEXT" json:"tags"`       // JSON-encoded array of strings
	Categories      string    `gorm:"type:TEXT" json:"categories"` // JSON-encoded array of strings
	InformationText string    `json:"information_text"`
}

var (
	ErrEmptyTag    = errors.New("empty tag")
	ErrTagNotFound = errors.New("tag not found")
)

func (s *Sheet) SaveSheet(db *gorm.DB) (*Sheet, error) {
	err := db.Model(&Sheet{}).Create(&s).Error
	if err != nil {
		return &Sheet{}, err
	}
	return s, nil
}

/*
log.Fatal(e) va planter le programme si le fichier n’existe pas, ce qui arrive pour les thumbnails manquants.
Pour corriger, il faut juste ignorer les fichiers non existants et loguer les erreurs autrement.
*/
func (s *Sheet) DeleteSheet(db *gorm.DB, sheetName string) (int64, error) {
	sheet, err := s.FindSheetBySafeName(db, sheetName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("Sheet not found")
		}
		return 0, err
	}

	paths := []string{
		path.Join(config.Config().ConfigPath, "sheets/uploaded-sheets", sheet.SafeComposer, sheet.SafeSheetName+".pdf"),
		path.Join(config.Config().ConfigPath, "sheets/thumbnails", sheet.SafeSheetName+".png"),
	}

	for _, filePath := range paths {
		err := os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Erreur lors de la suppression du fichier %s : %v\n", filePath, err)
		}
	}

	if sheet.SafeComposer == "unknown" {
		CheckAndDeleteUnknownComposer(db)
	}

	db = db.Model(&Sheet{}).Where("safe_sheet_name = ?", sheetName).Take(&Sheet{}).Delete(&Sheet{})

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("Sheet not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (s *Sheet) GetAllSheets(db *gorm.DB) (*[]Sheet, error) {
	/*
		This method will return max 20 sheets, to find more or specific one you need to specify it.
		Currently it sorts it by the newest updates
	*/
	var err error
	sheets := []Sheet{}

	err = db.Order("updated_at desc").Limit(20).Find(&sheets).Error
	if err != nil {
		return &[]Sheet{}, err
	}
	return &sheets, err
}

func (s *Sheet) FindSheetBySafeName(db *gorm.DB, sheetName string) (*Sheet, error) {
	// Get information of one single sheet by the safe sheet name
	var err error
	err = db.Model(&Sheet{}).Where("safe_sheet_name = ?", sheetName).Take(&s).Error
	if err != nil {
		return &Sheet{}, err
	}
	return s, nil
}

func (s *Sheet) List(db *gorm.DB, pagination Pagination, composer string) (*Pagination, error) {
	// For pagination

	var sheets []*Sheet
	if composer != "" {
		db.Scopes(ComposerEqual(composer), paginate(sheets, &pagination, db)).Find(&sheets)
	} else {
		db.Scopes(paginate(sheets, &pagination, db)).Find(&sheets)
	}

	pagination.Rows = sheets

	return &pagination, nil
}

func SearchSheet(db *gorm.DB, searchValue string) []*Sheet {
	// Search for sheets with containing string
	var sheets []*Sheet
	searchValue = "%" + searchValue + "%"
	db.Where("sheet_name LIKE ?", searchValue).Find(&sheets)
	return sheets
}

func ComposerEqual(composer string) func(db *gorm.DB) *gorm.DB {
	// Scope that composer is equal to composer (if you only want sheets from a certain composer)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("safe_composer = ?", composer)
	}
}

// Modification
// Méthode attachée à Sheet
// remplacement de s.Tags qui était de type []string (alias pq.StringArray)
// Solution limité à PostgreSQL
// Modification à string pour garder la compatibilité avec d'autres bases de données

// AppendTag ajoute une nouvelle tag à la feuille de musique, en évitant les doublons
func (s *Sheet) AppendTag(db *gorm.DB, tagValue string) error {
	tagValue = strings.TrimSpace(tagValue)
	if tagValue == "" {
		return ErrEmptyTag
	}

	var tags []string
	// Si s.Tags n'est pas vide, on essaie de le parser en JSON pour obtenir le slice de tags existants
	// Si s.Tags est vide, on part avec un slice de tags vide
	if s.Tags != "" {
		if err := json.Unmarshal([]byte(s.Tags), &tags); err != nil {
			return err
		}
	}

	// Vérifier si la tag existe déjà pour éviter les doublons
	for _, t := range tags {
		if t == tagValue {
			return nil
		}
	}

	tags = append(tags, tagValue)

	// Convertir le slice de tags en JSON
	// et stocker le résultat dans s.Tags
	// ["Classical","Romantic"]
	data, _ := json.Marshal(tags)
	s.Tags = string(data)

	return db.Save(s).Error
}

func (s *Sheet) DeleteTag(db *gorm.DB, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return ErrEmptyTag
	}

	var tags []string

	// 1️⃣ Décodage JSON
	if s.Tags != "" {
		if err := json.Unmarshal([]byte(s.Tags), &tags); err != nil {
			return err
		}
	}

	// 2️⃣ Recherche du tag
	index := -1
	for i, t := range tags {
		if t == value {
			index = i
			break
		}
	}

	if index == -1 {
		return ErrTagNotFound
	}

	// 3️⃣ Suppression
	tags = append(tags[:index], tags[index+1:]...)

	// 4️⃣ Ré-encodage
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	s.Tags = string(data)

	// 5️⃣ Sauvegarde
	return db.Save(s).Error
}

func (S *Sheet) UpdateSheetInformationText(db *gorm.DB, value string, sheet *Sheet) *Sheet {
	sheet.InformationText = value
	db.Save(sheet)

	return sheet
}

func FindSheetByTag(db *gorm.DB, tag string) ([]*Sheet, error) {
	tag = strings.TrimSpace(tag)

	// Validation de la tag
	if tag == "" {
		return nil, ErrEmptyTag
	}

	var allSheets []*Sheet
	var affectedSheets []*Sheet

	// Récupérer toutes les feuilles de musique
	if err := db.Find(&allSheets).Error; err != nil {
		return nil, err
	}

	// Parcourir toutes les feuilles et vérifier si la tag est présente dans le champ Tags de chaque feuille
	// Tags est stocké en JSON, il faut donc le décoder pour obtenir le slice de tags
	for _, sheet := range allSheets {

		var tags []string

		// 1️⃣ décoder JSON
		if sheet.Tags != "" {
			if err := json.Unmarshal([]byte(sheet.Tags), &tags); err != nil {
				continue // ignore sheet corrompue
			}
		}

		// 2️⃣ vérifier présence
		// Si la tag est présente, ajouter la feuille à la liste des feuilles affectées
		for _, t := range tags {
			if t == tag {
				affectedSheets = append(affectedSheets, sheet)
				break
			}
		}
	}

	return affectedSheets, nil
}
