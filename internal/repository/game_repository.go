package repository

import (
	"prismacrawler/internal/models"

	"gorm.io/gorm"
)

// GameRepositoryInterface define los métodos para interactuar con los datos de las partidas.
type GameRepositoryInterface interface {
	FindRunByIDWithCharacter(id uint) (*models.GameRun, error)
	SaveRunAndCharacter(run *models.GameRun, character *models.Character) error
	FindCharacterByIDAndUser(charID uint, userID uint) (*models.Character, error)
	FindMapByID(mapID uint) (*models.Map, error)
	CreateRun(run *models.GameRun) error
	GetTopRuns(limit int) ([]models.GameRun, error)
}

type gameRepository struct {
	db *gorm.DB
}

// NewGameRepository crea una nueva instancia del repositorio de partidas.
func NewGameRepository(db *gorm.DB) GameRepositoryInterface {
	return &gameRepository{db: db}
}

func (r *gameRepository) FindRunByIDWithCharacter(id uint) (*models.GameRun, error) {
	var run models.GameRun
	err := r.db.Preload("Character").First(&run, id).Error
	return &run, err
}

func (r *gameRepository) SaveRunAndCharacter(run *models.GameRun, character *models.Character) error {
	// Usamos una transacción para asegurar que ambos se guarden o ninguno.
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(run).Error; err != nil {
			return err
		}
		if err := tx.Save(character).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *gameRepository) FindCharacterByIDAndUser(charID uint, userID uint) (*models.Character, error) {
	var character models.Character
	err := r.db.Where("id = ? AND user_id = ?", charID, userID).First(&character).Error
	return &character, err
}

func (r *gameRepository) FindMapByID(mapID uint) (*models.Map, error) {
	var m models.Map
	err := r.db.First(&m, mapID).Error
	return &m, err
}

func (r *gameRepository) CreateRun(run *models.GameRun) error {
	return r.db.Create(run).Error
}

func (r *gameRepository) GetTopRuns(limit int) ([]models.GameRun, error) {
	var runs []models.GameRun
	err := r.db.Preload("Character").
		Where("status IN ?", []string{"Dead", "Won"}).
		Order("score desc, current_floor desc, kills desc").
		Limit(limit).
		Find(&runs).Error
	return runs, err
}
