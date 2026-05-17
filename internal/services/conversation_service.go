package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
	"github.com/google/uuid"
)

// ConversationService handles conversation business logic
type ConversationService struct {
	repo       *repository.ConversationRepository
	llmAPIKey  string
	llmBaseURL string
}

// NewConversationService creates a new conversation service
func NewConversationService(repo *repository.ConversationRepository) *ConversationService {
	return &ConversationService{
		repo:       repo,
		llmAPIKey:  os.Getenv("LLM_API_KEY"),
		llmBaseURL: os.Getenv("LLM_BASE_URL"),
	}
}

// StartAIChat starts a new AI conversation session
func (s *ConversationService) StartAIChat(userID string, scenarioID string, level int) (*models.StartChatResponse, error) {
	// Get the scenario
	scenario, err := s.repo.GetScenario(scenarioID)
	if err != nil {
		return nil, fmt.Errorf("scenario not found: %w", err)
	}

	// Check if user has an active session
	activeSession, err := s.repo.GetActiveSessionByUser(userID)
	if err != nil {
		return nil, err
	}

	// End any existing active session
	if activeSession != nil {
		s.repo.UpdateSessionStatus(activeSession.ID, "abandoned", nil)
	}

	// Create new session
	session := &models.ConversationSession{
		ID:         uuid.New().String(),
		UserID:     userID,
		Mode:       "ai",
		ScenarioID: scenarioID,
		Level:      level,
		Status:     "active",
		StartedAt:  time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Get starter messages from scenario
	var starterMessages []models.ConversationMessage
	
	// Add system message
	systemMsg := &models.ConversationMessage{
		ID:        uuid.New().String(),
		SessionID: session.ID,
		Sender:    "system",
		Content:   fmt.Sprintf("Scenario: %s - %s", scenario.Name, scenario.Description),
		CreatedAt: time.Now(),
	}
	s.repo.CreateMessage(systemMsg)
	starterMessages = append(starterMessages, *systemMsg)

	// Add AI starter message if available
	if scenario.StarterMessages != nil {
		if messages, ok := scenario.StarterMessages["messages"].([]interface{}); ok && len(messages) > 0 {
			// Pick first starter message
			starterContent := messages[0].(string)
			aiMsg := &models.ConversationMessage{
				ID:        uuid.New().String(),
				SessionID: session.ID,
				Sender:    "ai",
				Content:   starterContent,
				CreatedAt: time.Now(),
			}
			s.repo.CreateMessage(aiMsg)
			starterMessages = append(starterMessages, *aiMsg)
		}
	}

	return &models.StartChatResponse{
		Session:  *session,
		Messages: starterMessages,
		Scenario: *scenario,
	}, nil
}

// SendMessage processes a user message and gets AI response
func (s *ConversationService) SendMessage(sessionID string, userID string, message string) (*models.ChatResponse, error) {
	// Get session
	session, err := s.repo.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Verify session belongs to user
	if session.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to session")
	}

	// Get scenario
	scenario, err := s.repo.GetScenario(session.ScenarioID)
	if err != nil {
		return nil, fmt.Errorf("scenario not found: %w", err)
	}

	// Get conversation history
	history, err := s.repo.GetMessagesBySession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	// Save user message
	userMsg := &models.ConversationMessage{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Sender:    "user",
		Content:   message,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateMessage(userMsg); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Analyze user message for corrections
	correction, naturalnessScore := s.analyzeMessage(message, session.Level)

	// Update user message with correction info
	userMsg.Correction = correction
	userMsg.NaturalnessScore = naturalnessScore

	// Generate AI response
	aiResponse, err := s.generateAIResponse(message, history, scenario, session.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Save AI message
	aiMsg := &models.ConversationMessage{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Sender:    "ai",
		Content:   aiResponse,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateMessage(aiMsg); err != nil {
		return nil, fmt.Errorf("failed to save AI message: %w", err)
	}

	// Generate alternative phrasings
	alternatives := s.generateAlternatives(message, session.Level)

	return &models.ChatResponse{
		Message:          *userMsg,
		AIResponse:       *aiMsg,
		NaturalnessScore: naturalnessScore,
		Suggestions:      alternatives,
	}, nil
}

// analyzeMessage checks user input for grammar and naturalness
func (s *ConversationService) analyzeMessage(message string, level int) (string, int) {
	// This is a simplified analysis - in production, use proper NLP
	// For now, provide basic feedback based on message characteristics
	
	naturalnessScore := 70 // Default score
	correction := ""

	// Check for very short messages
	if len(message) < 5 {
		naturalnessScore -= 20
		correction = "もう少し詳しく話してみましょう。 (Try to be a bit more detailed.)"
	}

	// Check for casual vs polite forms based on level
	if level <= 3 {
		// N5-N3 should use polite forms
		if !containsPoliteForm(message) {
			naturalnessScore -= 15
			if correction != "" {
				correction += "\n"
			}
			correction += "です/ます形を使うとより丁寧ですよ。 (Using desu/masu forms would be more polite.)"
		}
	}

	// Check for natural flow markers
	if containsNaturalFillers(message) {
		naturalnessScore += 10
	}

	// Cap score
	if naturalnessScore > 100 {
		naturalnessScore = 100
	}
	if naturalnessScore < 0 {
		naturalnessScore = 0
	}

	if correction == "" && naturalnessScore >= 80 {
		correction = "いいね！自然な日本語です。 (Good! Natural Japanese.)"
	}

	return correction, naturalnessScore
}

// generateAIResponse creates an AI response based on context
func (s *ConversationService) generateAIResponse(userMessage string, history []models.ConversationMessage, scenario *models.ConversationScenario, level int) (string, error) {
	// If LLM API is configured, use it
	if s.llmAPIKey != "" {
		return s.callLLM(userMessage, history, scenario, level)
	}

	// Fallback: generate contextual response based on scenario
	return s.generateFallbackResponse(userMessage, scenario, level)
}

// callLLM makes an API call to the LLM service
func (s *ConversationService) callLLM(userMessage string, history []models.ConversationMessage, scenario *models.ConversationScenario, level int) (string, error) {
	// Build conversation context
	context := s.buildConversationContext(history)

	// Construct prompt
	prompt := fmt.Sprintf(`You are in this scenario: %s

%s

The user is at JLPT N%d level. Respond naturally in Japanese at an appropriate level for their proficiency.
Keep responses short (1-2 sentences) and conversational.

Conversation so far:
%s

User: %s

Respond as the AI character:`, scenario.Name, scenario.SystemPrompt, level, context, userMessage)

	// Make API call (OpenAI format as example)
	payload := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": scenario.SystemPrompt},
			{"role": "user", "content": prompt},
		},
		"max_tokens": 150,
		"temperature": 0.7,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.llmBaseURL+"/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.llmAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Extract response text
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("failed to parse LLM response")
}

// buildConversationContext creates a summary of recent conversation
func (s *ConversationService) buildConversationContext(history []models.ConversationMessage) string {
	// Get last 5 messages for context
	start := 0
	if len(history) > 5 {
		start = len(history) - 5
	}

	var context string
	for _, msg := range history[start:] {
		if msg.Sender == "user" {
			context += fmt.Sprintf("User: %s\n", msg.Content)
		} else if msg.Sender == "ai" {
			context += fmt.Sprintf("AI: %s\n", msg.Content)
		}
	}
	return context
}

// generateFallbackResponse creates a simple response when LLM is unavailable
func (s *ConversationService) generateFallbackResponse(userMessage string, scenario *models.ConversationScenario, level int) (string, error) {
	// Simple contextual responses based on scenario
	switch scenario.ID {
	case "konbini":
		return s.konbiniResponse(userMessage, level), nil
	case "restaurant":
		return s.restaurantResponse(userMessage, level), nil
	case "directions":
		return s.directionsResponse(userMessage, level), nil
	case "weather_chat":
		return s.weatherResponse(userMessage, level), nil
	case "weekend_plans":
		return s.weekendResponse(userMessage, level), nil
	default:
		return "そうですね。 (I see.)", nil
	}
}

// Scenario-specific response generators
func (s *ConversationService) konbiniResponse(userMessage string, level int) string {
	responses := []string{
		"はい、かしこまりました。",
		"お弁当、温めますか？",
		"ポイントカード、お持ちですか？",
		"袋、お分けしますか？",
		"ありがとうございました。",
	}
	// Simple rotation for demo
	return responses[time.Now().Second()%len(responses)]
}

func (s *ConversationService) restaurantResponse(userMessage string, level int) string {
	if containsAny(userMessage, []string{"おすすめ", "お勧め", "何がいい", "何がおいしい"}) {
		return "ラーメンがおすすめですよ。とても人気があります。"
	}
	if containsAny(userMessage, []string{"注文", "ちゅうもん", "これ", "これにする"}) {
		return "かしこまりました。お飲み物はいかがですか？"
	}
	if containsAny(userMessage, []string{"会計", "お会計", "かいけい", "おかいけい", "チェック"}) {
		return "かしこまりました。少々お待ちください。"
	}
	return "はい、何かご用ですか？"
}

func (s *ConversationService) directionsResponse(userMessage string, level int) string {
	if containsAny(userMessage, []string{"駅", "えき", "station"}) {
		return "駅はあの建物の隣です。まっすぐ行って、右に曲がってください。"
	}
	if containsAny(userMessage, []string{"トイレ", "toire", "bathroom", "restroom"}) {
		return "トイレはあそこです。あの看板の下にありますよ。"
	}
	return "すみません、よくわかりません。地図で調べてみてください。"
}

func (s *ConversationService) weatherResponse(userMessage string, level int) string {
	if containsAny(userMessage, []string{"暑い", "あつい", "hot"}) {
		return "そうですね。今日は本当に暑いですね。"
	}
	if containsAny(userMessage, []string{"寒い", "さむい", "cold"}) {
		return "そうですね。寒いですね。コート、着てきましたか？"
	}
	if containsAny(userMessage, []string{"雨", "あめ", "rain"}) {
		return "あ、ほんとだ。傘、持ってきましたか？"
	}
	return "そうですね。天気、気になりますね。"
}

func (s *ConversationService) weekendResponse(userMessage string, level int) string {
	if containsAny(userMessage, []string{"映画", "えいが", "movie"}) {
		return "いいね！何の映画を見るの？"
	}
	if containsAny(userMessage, []string{"友達", "ともだち", "友だち", "friend"}) {
		return "楽しみそう！どこに行くの？"
	}
	if containsAny(userMessage, []string{"家", "いえ", "うち", "home", "寝る", "ねる"}) {
		return "はは、休むのも大事だよね。"
	}
	return "へえ、いいね。楽しんでね！"
}

// generateAlternatives creates alternative phrasings for user messages
func (s *ConversationService) generateAlternatives(message string, level int) []string {
	var alternatives []string

	// Simple pattern-based alternatives
	if level <= 3 {
		// Suggest more polite forms
		if !containsPoliteForm(message) {
			alternatives = append(alternatives, "より丁寧に: 〜です/〜ます形を使う")
		}
	} else {
		// Suggest more casual/natural forms for higher levels
		if containsPoliteForm(message) {
			alternatives = append(alternatives, "よりカジュアルに: 〜だ/〜る形もOK")
		}
	}

	// Add filler word suggestions for naturalness
	if !containsNaturalFillers(message) {
		alternatives = append(alternatives, "自然な間: 「えーと」「なんか」を使ってみる")
	}

	return alternatives
}

// Helper functions
func containsPoliteForm(message string) bool {
	politeMarkers := []string{"です", "ます", "でした", "ました", "でしょう", "ましょう"}
	return containsAny(message, politeMarkers)
}

func containsNaturalFillers(message string) bool {
	fillers := []string{"えーと", "えっと", "なんか", "まぁ", "そうだな", "うーん"}
	return containsAny(message, fillers)
}

func containsAny(message string, substrings []string) bool {
	for _, s := range substrings {
		if len(s) <= len(message) {
			// Simple substring check
			for i := 0; i <= len(message)-len(s); i++ {
				if message[i:i+len(s)] == s {
					return true
				}
			}
		}
	}
	return false
}

// EndSession ends a conversation session
func (s *ConversationService) EndSession(sessionID string, userID string) error {
	// Get session
	session, err := s.repo.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Verify ownership
	if session.UserID != userID {
		return fmt.Errorf("unauthorized access to session")
	}

	// Calculate average naturalness from messages
	messages, err := s.repo.GetMessagesBySession(sessionID)
	if err != nil {
		return err
	}

	var totalScore int
	var scoreCount int
	for _, msg := range messages {
		if msg.Sender == "user" && msg.NaturalnessScore > 0 {
			totalScore += msg.NaturalnessScore
			scoreCount++
		}
	}

	var avgScore *int
	if scoreCount > 0 {
		avg := totalScore / scoreCount
		avgScore = &avg
	}

	return s.repo.UpdateSessionStatus(sessionID, "completed", avgScore)
}

// GetSessionHistory gets the full conversation history for a session
func (s *ConversationService) GetSessionHistory(sessionID string, userID string) (*models.ConversationSession, []models.ConversationMessage, error) {
	// Get session
	session, err := s.repo.GetSession(sessionID)
	if err != nil {
		return nil, nil, fmt.Errorf("session not found: %w", err)
	}

	// Verify ownership
	if session.UserID != userID {
		return nil, nil, fmt.Errorf("unauthorized access to session")
	}

	// Get messages
	messages, err := s.repo.GetMessagesBySession(sessionID)
	if err != nil {
		return nil, nil, err
	}

	return session, messages, nil
}

// GetUserStats gets conversation statistics for a user
func (s *ConversationService) GetUserStats(userID string) (map[string]interface{}, error) {
	return s.repo.GetSessionStats(userID)
}

// GetScenariosByLevel retrieves scenarios for a level
func (s *ConversationService) GetScenariosByLevel(level int, category string) ([]models.ConversationScenario, error) {
	// First seed default scenarios if needed
	s.repo.SeedScenarios()
	return s.repo.GetScenariosByLevel(level, category)
}
