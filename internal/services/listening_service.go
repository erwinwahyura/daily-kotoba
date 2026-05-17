package services

import (
	"fmt"

	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type ListeningService struct {
	listeningRepo *repository.ListeningRepository
	ttsService    *TTSService
}

func NewListeningService(listeningRepo *repository.ListeningRepository, ttsService *TTSService) *ListeningService {
	return &ListeningService{listeningRepo: listeningRepo, ttsService: ttsService}
}

func (s *ListeningService) GetExercises(level string) ([]models.ListeningExerciseSummary, error) {
	exercises, err := s.listeningRepo.GetExercisesByLevel(level)
	if err != nil {
		return nil, err
	}
	if exercises == nil {
		exercises = []models.ListeningExerciseSummary{}
	}
	return exercises, nil
}

func (s *ListeningService) GetExerciseDetail(id string) (*models.ListeningExerciseDetail, error) {
	ex, err := s.listeningRepo.GetExercise(id)
	if err != nil || ex == nil {
		return nil, fmt.Errorf("exercise not found")
	}

	questions, err := s.listeningRepo.GetQuestions(id)
	if err != nil {
		return nil, err
	}

	vocab, err := s.listeningRepo.GetVocabulary(id)
	if err != nil {
		return nil, err
	}

	// Lazily generate TTS audio if not cached
	audioURL := ""
	if ex.TTSCacheID != nil && *ex.TTSCacheID != "" {
		audioURL = fmt.Sprintf("/api/tts/audio/%s", *ex.TTSCacheID)
	} else if s.ttsService != nil && ex.Transcript != "" {
		ttsResp, ttsErr := s.ttsService.GenerateAudio(ex.Transcript, "")
		if ttsErr == nil && ttsResp != nil {
			audioURL = ttsResp.AudioURL
			// Extract audio ID from URL and cache it
			if len(ttsResp.AudioURL) > 15 {
				cacheID := ttsResp.AudioURL[len("/api/tts/audio/"):]
				_ = s.listeningRepo.SetTTSCacheID(id, cacheID)
			}
		}
	}

	return &models.ListeningExerciseDetail{
		ListeningExercise: *ex,
		Questions:         questions,
		Vocabulary:        vocab,
		AudioURL:          audioURL,
	}, nil
}

func (s *ListeningService) StartSession(userID, exerciseID string) (*models.ListeningSession, error) {
	questions, err := s.listeningRepo.GetQuestions(exerciseID)
	if err != nil {
		return nil, err
	}
	return s.listeningRepo.CreateSession(userID, exerciseID, len(questions))
}

func (s *ListeningService) SubmitAnswer(sessionID, questionID string, answerIdx int) (*models.SubmitListeningAnswerResponse, error) {
	q, err := s.listeningRepo.GetQuestion(questionID)
	if err != nil || q == nil {
		return nil, fmt.Errorf("question not found")
	}

	isCorrect, err := s.listeningRepo.SubmitAnswer(sessionID, questionID, answerIdx, q.CorrectIdx)
	if err != nil {
		return nil, err
	}

	return &models.SubmitListeningAnswerResponse{
		Correct:       isCorrect,
		CorrectAnswer: q.CorrectIdx,
	}, nil
}
