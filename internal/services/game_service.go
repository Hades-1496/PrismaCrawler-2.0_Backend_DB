package services

import (
	"context"
	"errors"
	"log"
	"time"

	"prismacrawler/internal/models"
	"prismacrawler/internal/repository"
	"prismacrawler/pkg/aiclient"
	"prismacrawler/pkg/utils"
)

var (
	ErrRunNotFound       = errors.New("partida no encontrada")
	ErrPermissionDenied  = errors.New("no tienes permiso para modificar esta partida")
	ErrRunAlreadyEnded   = errors.New("la partida ya ha finalizado")
	ErrCharacterNotFound = errors.New("personaje no encontrado o no te pertenece")
	ErrCharacterDead     = errors.New("este personaje está muerto y no puede iniciar una partida")
	ErrInternalCreate    = errors.New("error interno al crear la partida")
)

type StartRunInput struct {
	CharacterID uint
	MapID       uint
}

type SaveRunInput struct {
	RunID          uint
	CurrentFloor   int
	Score          int
	CurrentHP      int
	Kills          int
	DamageDealt    int
	DamageTaken    int
	ItemsCollected []string
}

type GameServiceInterface interface {
	StartRun(ctx context.Context, userID uint, input StartRunInput) (*models.GameRun, error)
	GetLeaderboard(ctx context.Context) ([]models.GameRun, error)
	SaveRun(ctx context.Context, userID uint, input SaveRunInput) (*models.GameRun, error)
}

// --- OCP: PATRÓN OBSERVER ---

// GameObserver permite extender la lógica de eventos de partidas sin modificar el servicio base
type GameObserver interface {
	OnSaveRun(run *models.GameRun, currentHP int, hpDropped, floorChanged bool)
}

// AIGameObserver implementa GameObserver para notificar al microservicio de Python
type AIGameObserver struct {
	ai *aiclient.Client
}

func NewAIGameObserver(ai *aiclient.Client) *AIGameObserver {
	return &AIGameObserver{ai: ai}
}

func (o *AIGameObserver) OnSaveRun(run *models.GameRun, currentHP int, hpDropped, floorChanged bool) {
	if o.ai == nil {
		return
	}
	if currentHP <= 0 {
		go func(runID uint, username, charName, status string, score, floor, damageDealt, damageTaken int) {
			ctxBg, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			payload := map[string]any{
				"event": "game_run_ended",
				"payload": map[string]any{
					"run_id":       runID,
					"username":     username,
					"character":    charName,
					"status":       status,
					"score":        score,
					"floor":        floor,
					"damage_dealt": damageDealt,
					"damage_taken": damageTaken,
				},
			}
			if _, err := o.ai.Proxy(ctxBg, "/api/n8n/relay", payload); err != nil {
				log.Printf("Error enviando trigger game_run_ended a IA: %v", err)
			}
		}(run.ID, run.Character.User.Nickname, run.Character.Name, run.Status, run.Score, run.CurrentFloor, run.DamageDealt, run.DamageTaken)
	} else if floorChanged || hpDropped {
		go func(runID uint, charName string, floor, hp int) {
			ctxBg, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			eventType := "floor_advanced"
			if hpDropped {
				eventType = "massive_damage_taken"
			}
			payload := map[string]any{
				"event": eventType,
				"data": map[string]any{
					"run_id":         runID,
					"character_name": charName,
					"floor":          floor,
					"current_hp":     hp,
				},
			}
			o.ai.Proxy(ctxBg, "/api/game/event", payload)
		}(run.ID, run.Character.Name, run.CurrentFloor, run.Character.BaseHP)
	}
}

// --- GAME SERVICE ---

type gameService struct {
	repo      repository.GameRepositoryInterface
	observers []GameObserver // Aceptamos múltiples observadores inyectados
}

// NewGameService acepta un número variable de observadores (variadic)
func NewGameService(repo repository.GameRepositoryInterface, observers ...GameObserver) GameServiceInterface {
	return &gameService{repo: repo, observers: observers}
}

func (s *gameService) StartRun(ctx context.Context, userID uint, input StartRunInput) (*models.GameRun, error) {
	character, err := s.repo.FindCharacterByIDAndUser(input.CharacterID, userID)
	if err != nil {
		return nil, ErrCharacterNotFound
	}
	if !character.IsAlive {
		return nil, ErrCharacterDead
	}

	run := models.GameRun{
		CharacterID: character.ID,
	}

	if input.MapID > 0 {
		if _, err := s.repo.FindMapByID(input.MapID); err == nil {
			mapID := input.MapID
			run.MapID = &mapID
		}
	} else {
		run.Seed = utils.GenerateSeed(6)
	}

	if err := s.repo.CreateRun(&run); err != nil {
		return nil, ErrInternalCreate
	}

	return &run, nil
}

func (s *gameService) GetLeaderboard(ctx context.Context) ([]models.GameRun, error) {
	return s.repo.GetTopRuns(10)
}

func (s *gameService) SaveRun(ctx context.Context, userID uint, input SaveRunInput) (*models.GameRun, error) {
	run, err := s.repo.FindRunByIDWithCharacter(input.RunID)
	if err != nil {
		return nil, ErrRunNotFound
	}

	if !run.IsOwnedBy(userID) {
		return nil, ErrPermissionDenied
	}
	if run.Status != "In_Progress" {
		return nil, ErrRunAlreadyEnded
	}

	// Detectar hitos para el Game Observer antes de sobreescribir el estado
	hpDropped := input.CurrentHP > 0 && run.Character.BaseHP-input.CurrentHP >= 30
	floorChanged := run.CurrentFloor < input.CurrentFloor

	run.UpdateState(input.CurrentFloor, input.Score, input.CurrentHP, input.Kills, input.DamageDealt, input.DamageTaken)

	if err := s.repo.SaveRunAndCharacter(run, &run.Character); err != nil {
		return nil, err
	}

	if len(input.ItemsCollected) > 0 {
		_ = s.repo.AddItemsToRun(run.ID, input.ItemsCollected)
	}

	// SRP / OCP: Notificamos a los observadores registrados dinámicamente
	for _, obs := range s.observers {
		if obs != nil {
			obs.OnSaveRun(run, input.CurrentHP, hpDropped, floorChanged)
		}
	}

	return run, nil
}
