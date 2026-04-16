package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type ConjugationService struct {
	conjRepo *repository.ConjugationRepository
}

func NewConjugationService(conjRepo *repository.ConjugationRepository) *ConjugationService {
	return &ConjugationService{conjRepo: conjRepo}
}

// StartDrillSession starts a new conjugation drill session for a user
func (s *ConjugationService) StartDrillSession(userID string, targetForm string) (*models.ConjugationSessionResponse, error) {
	// Get user's current level (simplified - could fetch from user profile)
	jlptLevel := "N4" // Default, could be determined from user progress

	// Get challenges for the form
	challenges, err := s.conjRepo.GetChallengesByForm(targetForm, jlptLevel, 10)
	if err != nil {
		return nil, err
	}
	if len(challenges) == 0 {
		return nil, fmt.Errorf("no challenges found for form: %s", targetForm)
	}

	// Create session
	session := &models.ConjugationSession{
		ID:             uuid.New().String(),
		UserID:         userID,
		CurrentForm:    targetForm,
		CurrentIndex:   0,
		TotalQuestions: len(challenges),
		StartTime:      time.Now(),
		LastActive:     time.Now(),
	}

	if err := s.conjRepo.CreateSession(session); err != nil {
		return nil, err
	}

	// Get form info
	formInfo := s.getFormInfo(targetForm)

	return &models.ConjugationSessionResponse{
		Session:    session,
		Challenges: challenges,
		Progress: &models.ConjugationProgress{
			CurrentForm:      targetForm,
			TotalAttempts:    0,
			CurrentStreak:    0,
			DailyGoal:        20,
			DailyCompleted:   0,
		},
		FormInfo: formInfo,
	}, nil
}

// GetNextChallenge gets the next challenge in a session
func (s *ConjugationService) GetNextChallenge(sessionID, userID string) (*models.ConjugationChallengeResponse, error) {
	// Get session
	session, err := s.conjRepo.GetActiveSession(userID)
	if err != nil || session == nil || session.ID != sessionID {
		return nil, fmt.Errorf("session not found or expired")
	}

	// Get challenges for current form
	challenges, err := s.conjRepo.GetChallengesByForm(session.CurrentForm, "N4", 20)
	if err != nil {
		return nil, err
	}

	if session.CurrentIndex >= len(challenges) {
		// Form completed, could advance to next form
		return nil, fmt.Errorf("form completed")
	}

	formInfo := s.getFormInfo(session.CurrentForm)

	return &models.ConjugationChallengeResponse{
		Challenge: challenges[session.CurrentIndex],
		Progress: &models.ConjugationProgress{
			CurrentForm:     session.CurrentForm,
			CurrentStreak:   session.Streak,
			BestStreak:      session.MaxStreak,
			DailyGoal:       20,
			DailyCompleted:  session.CorrectCount,
		},
		FormInfo:  formInfo,
		SessionID: session.ID,
	}, nil
}

// SubmitAnswer checks user's answer and returns result
func (s *ConjugationService) SubmitAnswer(userID, sessionID, challengeID, answer string, timeSpentMs int) (*models.ConjugationSubmitResponse, error) {
	// Get challenge
	challenge, err := s.conjRepo.GetChallengeByID(challengeID)
	if err != nil {
		return nil, err
	}

	// Get session
	session, err := s.conjRepo.GetActiveSession(userID)
	if err != nil || session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Check answer (flexible matching)
	isCorrect := s.checkAnswer(answer, challenge.FullAnswer, challenge.TargetEnding)

	// Update session
	session.TotalQuestions++
	if isCorrect {
		session.CorrectCount++
		session.Streak++
		if session.Streak > session.MaxStreak {
			session.MaxStreak = session.Streak
		}
	} else {
		session.WrongCount++
		session.Streak = 0
	}
	session.CurrentIndex++
	session.LastActive = time.Now()

	// Record attempt
	attempt := &models.ConjugationAttempt{
		ID:           uuid.New().String(),
		SessionID:    sessionID,
		UserID:       userID,
		ChallengeID:  challengeID,
		FormType:     challenge.TargetForm,
		BaseForm:     challenge.BaseForm,
		UserAnswer:   answer,
		IsCorrect:    isCorrect,
		TimeSpentSec: timeSpentMs / 1000,
		CreatedAt:    time.Now(),
	}

	if err := s.conjRepo.RecordAttempt(attempt); err != nil {
		// Log but don't fail
	}

	if err := s.conjRepo.UpdateSession(session); err != nil {
		return nil, err
	}

	// Get next challenge or mark form complete
	var nextChallenge *models.ConjugationChallenge
	var nextFormInfo *models.ConjugationFormType
	formCompleted := false
	allFormsCompleted := false

	challenges, err := s.conjRepo.GetChallengesByForm(session.CurrentForm, "N4", 20)
	if err == nil && session.CurrentIndex < len(challenges) {
		nextChallenge = challenges[session.CurrentIndex]
		nextFormInfo = s.getFormInfo(session.CurrentForm)
	} else {
		formCompleted = true
		// Could check if all forms completed
	}

	explanation := s.getExplanation(challenge, isCorrect)

	return &models.ConjugationSubmitResponse{
		IsCorrect:         isCorrect,
		CorrectAnswer:     challenge.FullAnswer,
		Explanation:       explanation,
		NextChallenge:     nextChallenge,
		NextFormInfo:      nextFormInfo,
		FormCompleted:     formCompleted,
		AllFormsCompleted: allFormsCompleted,
		Progress: &models.ConjugationProgress{
			CurrentForm:      session.CurrentForm,
			CurrentStreak:    session.Streak,
			BestStreak:       session.MaxStreak,
			AccuracyRate:     float64(session.CorrectCount) / float64(session.CorrectCount+session.WrongCount) * 100,
			DailyGoal:        20,
			DailyCompleted:   session.CorrectCount,
		},
		SessionID: session.ID,
	}, nil
}

// GetProgress retrieves user's conjugation progress
func (s *ConjugationService) GetProgress(userID string) (*models.ConjugationProgress, error) {
	return s.conjRepo.GetProgressStats(userID)
}

// Helper: check answer with flexible matching
func (s *ConjugationService) checkAnswer(userAnswer, correctAnswer, targetEnding string) bool {
	// Normalize
	user := []rune(userAnswer)
	correct := []rune(correctAnswer)

	// Check full match
	if string(user) == string(correct) {
		return true
	}

	// Check if user provided just the ending correctly
	if len(user) > 0 {
		// Extract ending from correct answer
		// For 食べて, ending is て; for 行かない, ending is ない
		endingOK := false
		if string(user) == targetEnding {
			endingOK = true
		}

		// Also accept full form
		if endingOK {
			return true
		}
	}

	// Allow small typos (one character difference)
	if len(user) == len(correct) {
		diff := 0
		for i := range user {
			if user[i] != correct[i] {
				diff++
			}
		}
		if diff <= 1 {
			return true
		}
	}

	return false
}

// Helper: get form info
func (s *ConjugationService) getFormInfo(formName string) *models.ConjugationFormType {
	forms := models.GetConjugationForms()
	for _, f := range forms {
		if f.Name == formName {
			return &f
		}
	}
	return nil
}

// Helper: generate explanation
func (s *ConjugationService) getExplanation(challenge *models.ConjugationChallenge, isCorrect bool) string {
	if isCorrect {
		return fmt.Sprintf("Correct! %s → %s (%s form)", challenge.BaseForm, challenge.FullAnswer, challenge.TargetForm)
	}
	return fmt.Sprintf("The correct answer is %s. Hint: %s", challenge.FullAnswer, challenge.Hint)
}

// GetWeakPointsAnalysis analyzes user's conjugation performance
func (s *ConjugationService) GetWeakPointsAnalysis(userID string) (*models.WeakPointsAnalysis, error) {
	// Get accuracy stats per form
	weakPoints, err := s.conjRepo.GetWeakPointsByForm(userID)
	if err != nil {
		return nil, err
	}
	
	// Identify weak forms (accuracy < 70%)
	var weakForms []models.WeakForm
	var strongForms []models.WeakForm
	
	for form, data := range weakPoints {
		accuracy := data["accuracy"].(float64)
		total := data["total"].(int)
		
		wf := models.WeakForm{
			Form:     form,
			Accuracy: accuracy,
			Total:    total,
		}
		
		if accuracy < 70.0 && total >= 5 {
			weakForms = append(weakForms, wf)
		} else if accuracy >= 80.0 {
			strongForms = append(strongForms, wf)
		}
	}
	
	// Sort weak forms by accuracy (lowest first)
	for i := 0; i < len(weakForms); i++ {
		for j := i + 1; j < len(weakForms); j++ {
			if weakForms[i].Accuracy > weakForms[j].Accuracy {
				weakForms[i], weakForms[j] = weakForms[j], weakForms[i]
			}
		}
	}
	
	return &models.WeakPointsAnalysis{
		WeakForms:   weakForms,
		StrongForms: strongForms,
		TotalForms:  len(weakPoints),
	}, nil
}

// GenerateWeakPointDrill creates a focused drill for weak forms
func (s *ConjugationService) GenerateWeakPointDrill(userID string) (*models.ConjugationSessionResponse, error) {
	// Get weak points analysis
	analysis, err := s.GetWeakPointsAnalysis(userID)
	if err != nil {
		return nil, err
	}
	
	// If no weak forms, return error
	if len(analysis.WeakForms) == 0 {
		return nil, fmt.Errorf("no weak points found - you're doing great!")
	}
	
	// Pick the weakest form
	targetForm := analysis.WeakForms[0].Form
	
	// Get challenges for this weak form (prioritizing ones user got wrong)
	challenges, err := s.conjRepo.GetChallengesForWeakPoint(targetForm, userID, 10)
	if err != nil {
		return nil, err
	}
	
	if len(challenges) == 0 {
		return nil, fmt.Errorf("no challenges available for form: %s", targetForm)
	}
	
	// Create session
	session := &models.ConjugationSession{
		ID:             uuid.New().String(),
		UserID:         userID,
		CurrentForm:    targetForm,
		CurrentIndex:   0,
		TotalQuestions: len(challenges),
		StartTime:      time.Now(),
		LastActive:     time.Now(),
		IsWeakPointDrill: true,
		TargetWeakForm: targetForm,
	}
	
	if err := s.conjRepo.CreateSession(session); err != nil {
		return nil, err
	}
	
	// Get form info
	formInfo := s.getFormInfo(targetForm)
	
	return &models.ConjugationSessionResponse{
		Session:    session,
		Challenges: challenges,
		Progress: &models.ConjugationProgress{
			CurrentForm:      targetForm,
			TotalAttempts:    0,
			CurrentStreak:    0,
			DailyGoal:        10,
			DailyCompleted:   0,
			IsWeakPointDrill: true,
			WeakFormAccuracy: analysis.WeakForms[0].Accuracy,
		},
		FormInfo: formInfo,
	}, nil
}
