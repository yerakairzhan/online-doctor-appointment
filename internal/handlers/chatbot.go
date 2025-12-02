package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"online-doctor-appointment/internal/database"
	"os"
	"strings"
)

// ChatRequest represents a chat message from the user
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse represents the AI response
type ChatResponse struct {
	Response  string `json:"response"`
	Timestamp string `json:"timestamp"`
}

// Gemini API structures
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// ChatbotPageHandler serves the chatbot page
func ChatbotPageHandler(w http.ResponseWriter, r *http.Request) {
	// You can add authentication check here if needed
	// For now, we'll make it accessible to everyone

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Medical Assistant - Online Doctor Appointment</title>
    <link rel="stylesheet" href="/static/css/style.css">
    <style>
        .chat-container {
            background: white;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            overflow: hidden;
            display: flex;
            flex-direction: column;
            height: 600px;
            max-width: 1000px;
            margin: 0 auto;
        }

        .chat-header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .chat-header h2 {
            font-size: 1.5rem;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .bot-avatar {
            width: 40px;
            height: 40px;
            background: white;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 1.5rem;
        }

        .chat-messages {
            flex: 1;
            padding: 20px;
            overflow-y: auto;
            background: #f8f9fa;
        }

        .message {
            display: flex;
            margin-bottom: 20px;
            animation: fadeIn 0.3s ease-in;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .message.bot {
            justify-content: flex-start;
        }

        .message.user {
            justify-content: flex-end;
        }

        .message-avatar {
            width: 36px;
            height: 36px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 1.2rem;
            flex-shrink: 0;
        }

        .message.bot .message-avatar {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            margin-right: 10px;
        }

        .message.user .message-avatar {
            background: #28a745;
            color: white;
            margin-left: 10px;
            order: 2;
        }

        .message-content {
            max-width: 70%;
            padding: 12px 16px;
            border-radius: 12px;
            word-wrap: break-word;
        }

        .message.bot .message-content {
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 12px 12px 12px 0;
        }

        .message.user .message-content {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border-radius: 12px 12px 0 12px;
        }

        .message-time {
            font-size: 0.75rem;
            color: #999;
            margin-top: 5px;
        }

        .typing-indicator {
            display: flex;
            align-items: center;
            gap: 5px;
            padding: 10px 16px;
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 12px 12px 12px 0;
            width: fit-content;
        }

        .typing-dot {
            width: 8px;
            height: 8px;
            background: #667eea;
            border-radius: 50%;
            animation: typing 1.4s infinite;
        }

        .typing-dot:nth-child(2) { animation-delay: 0.2s; }
        .typing-dot:nth-child(3) { animation-delay: 0.4s; }

        @keyframes typing {
            0%, 60%, 100% { transform: translateY(0); }
            30% { transform: translateY(-10px); }
        }

        .quick-questions {
            padding: 15px 20px;
            background: white;
            border-top: 1px solid #e9ecef;
        }

        .quick-questions-title {
            font-size: 0.9rem;
            color: #666;
            margin-bottom: 10px;
        }

        .quick-questions-list {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }

        .quick-question-btn {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            padding: 8px 16px;
            border-radius: 20px;
            cursor: pointer;
            transition: all 0.3s;
            font-size: 0.9rem;
            color: #667eea;
        }

        .quick-question-btn:hover {
            background: #667eea;
            color: white;
            transform: translateY(-2px);
        }

        .chat-input-area {
            padding: 20px;
            background: white;
            border-top: 1px solid #e9ecef;
        }

        .chat-input-container {
            display: flex;
            gap: 10px;
            align-items: center;
        }

        .chat-input {
            flex: 1;
            padding: 12px 16px;
            border: 2px solid #e9ecef;
            border-radius: 25px;
            font-size: 1rem;
            transition: border-color 0.3s;
        }

        .chat-input:focus {
            outline: none;
            border-color: #667eea;
        }

        .btn-send {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s;
        }

        .btn-send:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }

        .disclaimer {
            background: #fff3cd;
            border: 1px solid #ffc107;
            color: #856404;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 20px;
            font-size: 0.9rem;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🤖 AI Medical Assistant</h1>
            <p>Ask questions about symptoms, medications, and general health</p>
            <a href="/dashboard/patient" class="btn btn-secondary">← Back to Dashboard</a>
        </div>

        <div class="disclaimer">
            ⚠️ <strong>Disclaimer:</strong> This AI assistant provides general health information only. For medical emergencies, call emergency services. Always consult a qualified healthcare professional for diagnosis and treatment.
        </div>

        <div class="chat-container">
            <div class="chat-header">
                <h2>
                    <div class="bot-avatar">🤖</div>
                    MediBot Assistant
                </h2>
            </div>

            <div class="chat-messages" id="chatMessages">
                <div class="message bot">
                    <div class="message-avatar">🤖</div>
                    <div>
                        <div class="message-content">
                            Hello! I'm your AI Medical Assistant. I can help you with general health questions, information about symptoms, medications, and wellness tips. What would you like to know?
                        </div>
                        <div class="message-time">Just now</div>
                    </div>
                </div>
            </div>

            <div class="quick-questions">
                <div class="quick-questions-title">💡 Quick Questions:</div>
                <div class="quick-questions-list">
                    <button class="quick-question-btn" onclick="askQuestion('What are common cold symptoms?')">
                        Common cold symptoms?
                    </button>
                    <button class="quick-question-btn" onclick="askQuestion('How to lower blood pressure naturally?')">
                        Lower blood pressure?
                    </button>
                    <button class="quick-question-btn" onclick="askQuestion('What causes headaches?')">
                        What causes headaches?
                    </button>
                    <button class="quick-question-btn" onclick="askQuestion('How much water should I drink daily?')">
                        Daily water intake?
                    </button>
                </div>
            </div>

            <div class="chat-input-area">
                <div class="chat-input-container">
                    <input 
                        type="text" 
                        class="chat-input" 
                        id="chatInput"
                        placeholder="Type your health question here..."
                        onkeypress="handleKeyPress(event)"
                    >
                    <button class="btn-send" onclick="sendMessage()">
                        Send 📤
                    </button>
                </div>
            </div>
        </div>
    </div>

    <script>
        function getCurrentTime() {
            const now = new Date();
            return now.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
        }

        function addMessage(content, isUser) {
            const messagesContainer = document.getElementById('chatMessages');
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message ' + (isUser ? 'user' : 'bot');
            
            messageDiv.innerHTML = 
                '<div class="message-avatar">' + (isUser ? '👤' : '🤖') + '</div>' +
                '<div>' +
                    '<div class="message-content">' + content + '</div>' +
                    '<div class="message-time">' + getCurrentTime() + '</div>' +
                '</div>';
            
            messagesContainer.appendChild(messageDiv);
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }

        function showTypingIndicator() {
            const messagesContainer = document.getElementById('chatMessages');
            const typingDiv = document.createElement('div');
            typingDiv.className = 'message bot';
            typingDiv.id = 'typingIndicator';
            
            typingDiv.innerHTML = 
                '<div class="message-avatar">🤖</div>' +
                '<div class="typing-indicator">' +
                    '<div class="typing-dot"></div>' +
                    '<div class="typing-dot"></div>' +
                    '<div class="typing-dot"></div>' +
                '</div>';
            
            messagesContainer.appendChild(typingDiv);
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }

        function removeTypingIndicator() {
            const typingIndicator = document.getElementById('typingIndicator');
            if (typingIndicator) {
                typingIndicator.remove();
            }
        }

        async function sendMessage() {
            const input = document.getElementById('chatInput');
            const message = input.value.trim();
            
            if (message === '') return;
            
            addMessage(message, true);
            input.value = '';
            
            showTypingIndicator();
            
            try {
                const response = await fetch('/api/chatbot', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ message: message })
                });
                
                const data = await response.json();
                removeTypingIndicator();
                
                if (response.ok) {
                    addMessage(data.response, false);
                } else {
                    addMessage('Sorry, I encountered an error. Please try again.', false);
                }
            } catch (error) {
                removeTypingIndicator();
                addMessage('Sorry, I could not connect to the server. Please check your connection.', false);
            }
        }

        function askQuestion(question) {
            const input = document.getElementById('chatInput');
            input.value = question;
            sendMessage();
        }

        function handleKeyPress(event) {
            if (event.key === 'Enter') {
                sendMessage();
            }
        }

        window.onload = function() {
            document.getElementById('chatInput').focus();
        };
    </script>
</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// ChatbotAPIHandler handles chatbot API requests
func ChatbotAPIHandler(w http.ResponseWriter, r *http.Request) {
	var chatReq ChatRequest

	err := json.NewDecoder(r.Body).Decode(&chatReq)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	// Check if using local AI or API
	useLocalAI := os.Getenv("USE_LOCAL_AI") == "true"

	var response string
	if useLocalAI {
		response = getLocalAIResponse(chatReq.Message)
	} else {
		response, err = getGeminiResponse(chatReq.Message)
		if err != nil {
			fmt.Println("Gemini API error:", err)
			// Fallback to local responses if API fails
			response = getLocalAIResponse(chatReq.Message)
		}
	}

	chatResp := ChatResponse{
		Response:  response,
		Timestamp: getCurrentTimestamp(),
	}

	userID, _, _ := GetCurrentUser(r)
	go saveChatToDB(userID, chatReq.Message, response)

	respondWithJSON(w, http.StatusOK, chatResp)
}

// getGeminiResponse calls Google Gemini API
func getGeminiResponse(userMessage string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Gemini API key not configured")
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-2.0-flash-exp"
	}

	systemPrompt := `You are a helpful medical assistant chatbot for an online doctor appointment system. 
Provide accurate, helpful medical information while being careful to:
1. Always recommend consulting a healthcare professional for serious concerns
2. Never provide specific diagnoses
3. Give general health advice and information about symptoms
4. Be empathetic and professional
5. Keep responses concise (2-3 sentences)
6. Include a reminder to book an appointment for personalized care when appropriate`

	// Combine system prompt with user message for Gemini
	fullMessage := systemPrompt + "\n\nUser question: " + userMessage

	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: fullMessage},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Construct the API URL with the model and API key
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// getLocalAIResponse provides pre-programmed responses (fallback)
func getLocalAIResponse(userMessage string) string {
	lowerMessage := strings.ToLower(userMessage)

	responses := map[string]string{
		"cold": "Common cold symptoms include runny nose, sore throat, cough, congestion, mild body aches, sneezing, and low-grade fever. Most colds resolve within 7-10 days with rest and fluids. If symptoms persist or worsen, please book an appointment with one of our doctors for a proper evaluation.",

		"blood pressure": "To naturally lower blood pressure: 1) Exercise regularly (30 min/day), 2) Reduce sodium intake, 3) Eat potassium-rich foods, 4) Limit alcohol, 5) Manage stress, 6) Maintain healthy weight. Always monitor with your doctor. Would you like to book an appointment with our cardiologist?",

		"headache": "Common headache causes include tension, dehydration, lack of sleep, stress, eye strain, sinus issues, or caffeine withdrawal. Stay hydrated, rest, and manage stress. Frequent or severe headaches need medical evaluation. Book an appointment if headaches persist or worsen.",

		"water": "General recommendation is 8 glasses (2 liters) per day, but needs vary based on activity level, climate, and health conditions. A good indicator is pale yellow urine. Increase intake during exercise or hot weather.",

		"fever": "A fever (temperature above 38°C/100.4°F) is usually a sign your body is fighting an infection. Rest, drink fluids, and take fever reducers if needed. Seek immediate medical attention if fever exceeds 39.4°C (103°F) or lasts more than 3 days. Book an appointment with our doctors for proper evaluation.",

		"diabetes": "Diabetes is a condition affecting blood sugar regulation. Common symptoms include increased thirst, frequent urination, fatigue, and blurred vision. Management includes diet, exercise, medication, and regular monitoring. Our endocrinology specialists can provide personalized care - would you like to book an appointment?",

		"covid": "COVID-19 symptoms include fever, cough, fatigue, loss of taste/smell, and difficulty breathing. If you suspect COVID-19, get tested and isolate. Severe symptoms require immediate medical attention. Our doctors can provide telehealth consultations. Would you like to book an appointment?",
	}

	for key, response := range responses {
		if strings.Contains(lowerMessage, key) {
			return response
		}
	}

	// Default response
	return "That's an interesting question! For specific medical concerns and personalized advice, I recommend consulting with one of our qualified healthcare professionals. They can provide guidance based on your medical history and current health status. Would you like to book an appointment with one of our doctors? You can do so from your dashboard."
}

func getCurrentTimestamp() string {
	// Return current timestamp in ISO format
	return ""
}

func saveChatToDB(userID int, message, response string) {
	if database.DB == nil {
		fmt.Println("DB not initialized")
		return
	}

	_, err := database.DB.Exec(`
        INSERT INTO chat_logs (user_id, message, response)
        VALUES ($1, $2, $3)
    `, userID, message, response)

	if err != nil {
		fmt.Println("DB insert error:", err)
	}
}
