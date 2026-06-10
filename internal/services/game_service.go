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
	ErrRunNotFound      = errors.New("partida no encontrada")
	ErrPermissionDenied = errors.New("no tienes permiso para modificar esta partida")
	ErrRunAlreadyEnded  = errors.New("la partida ya ha finalizado")
	ErrCharacterNotFound = errors.New("personaje no encontrado o no te pertenece")
	ErrCharacterDead     = errors.New("este personaje está muerto y no puede iniciar una partida")
	ErrInternalCreate    = errors.New("error interno al crear la partida")
)

type StartRunInput struct {
	CharacterID uint
	MapID       uint
}

type SaveRunInput struct {
	RunID        uint
	CurrentFloor int
	Score        int
	CurrentHP    int
	Kills        int
	DamageDealt  int
	DamageTaken  int
}

type GameServiceInterface interface {
	StartRun(ctx context.Context, userID uint, input StartRunInput) (*models.GameRun, error)
	GetLeaderboard(ctx context.Context) ([]models.GameRun, error)
	SaveRun(ctx context.Context, userID uint, input SaveRunInput) (*models.GameRun, error)
}

type gameService struct {
	repo repository.GameRepositoryInterface
	ai   *aiclient.Client
}

func NewGameService(repo repository.GameRepositoryInterface, ai *aiclient.Client) GameServiceInterface {
	return &gameService{repo: repo, ai: ai}
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

	run.UpdateState(input.CurrentFloor, input.Score, input.CurrentHP, input.Kills, input.DamageDealt, input.DamageTaken)

	if err := s.repo.SaveRunAndCharacter(run, &run.Character); err != nil {
		return nil, err
	}

	if input.CurrentHP <= 0 && s.ai != nil {
		go func(runID uint, charName string, score, floor int) {
			ctxBg, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			payload := map[string]any{
				"event": "game_run_ended",
				"data": map[string]any{
					"run_id":    runID,
					"character": charName,
					"score":     score,
					"floor":     floor,
				},
			}
			if _, err := s.ai.Proxy(ctxBg, "/api/n8n/relay", payload); err != nil {
				log.Printf("Error enviando trigger game_run_ended a IA: %v", err)
			}
		}(run.ID, run.Character.Name, run.Score, run.CurrentFloor)
	}

	return run, nil
}
