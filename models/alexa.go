package models

// AlexaRequest representa a requisição enviada pelo serviço da Alexa
type AlexaRequest struct {
	Version string       `json:"version"`
	Session AlexaSession `json:"session"`
	Request AlexaPayload `json:"request"`
}

type AlexaSession struct {
	New         bool                   `json:"new"`
	SessionID   string                 `json:"sessionId"`
	Application AlexaApplication       `json:"application"`
	Attributes  map[string]interface{} `json:"attributes"`
	User        AlexaUser              `json:"user"`
}

type AlexaApplication struct {
	ApplicationID string `json:"applicationId"`
}

type AlexaUser struct {
	UserID string `json:"userId"`
}

type AlexaPayload struct {
	Type      string      `json:"type"` // LaunchRequest, IntentRequest, SessionEndedRequest
	RequestID string      `json:"requestId"`
	Timestamp string      `json:"timestamp"`
	Locale    string      `json:"locale"`
	Intent    AlexaIntent `json:"intent"`
	Reason    string      `json:"reason"`
}

type AlexaIntent struct {
	Name               string               `json:"name"`
	ConfirmationStatus string               `json:"confirmationStatus"`
	Slots              map[string]AlexaSlot `json:"slots"`
}

type AlexaSlot struct {
	Name               string `json:"name"`
	Value              string `json:"value"`
	ConfirmationStatus string `json:"confirmationStatus"`
}

// AlexaResponse representa a resposta devolvida para a Alexa
type AlexaResponse struct {
	Version           string            `json:"version"`
	SessionAttributes map[string]any    `json:"sessionAttributes,omitempty"`
	Response          AlexaResponseBody `json:"response"`
}

type AlexaResponseBody struct {
	OutputSpeech     *AlexaOutputSpeech `json:"outputSpeech,omitempty"`
	Reprompt         *AlexaReprompt     `json:"reprompt,omitempty"`
	ShouldEndSession bool               `json:"shouldEndSession"`
}

type AlexaOutputSpeech struct {
	Type string `json:"type"` // "PlainText" ou "SSML"
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

type AlexaReprompt struct {
	OutputSpeech AlexaOutputSpeech `json:"outputSpeech"`
}

// BuildAlexaResponse constrói uma resposta simples da Alexa
func BuildAlexaResponse(speechText string, endSession bool) AlexaResponse {
	return AlexaResponse{
		Version: "1.0",
		Response: AlexaResponseBody{
			OutputSpeech: &AlexaOutputSpeech{
				Type: "PlainText",
				Text: speechText,
			},
			ShouldEndSession: endSession,
		},
	}
}

// BuildAlexaRepromptResponse constrói uma resposta mantendo a sessão aberta para resposta
func BuildAlexaRepromptResponse(speechText, repromptText string) AlexaResponse {
	return AlexaResponse{
		Version: "1.0",
		Response: AlexaResponseBody{
			OutputSpeech: &AlexaOutputSpeech{
				Type: "PlainText",
				Text: speechText,
			},
			Reprompt: &AlexaReprompt{
				OutputSpeech: AlexaOutputSpeech{
					Type: "PlainText",
					Text: repromptText,
				},
			},
			ShouldEndSession: false,
		},
	}
}
