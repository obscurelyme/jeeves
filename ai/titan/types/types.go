package titan

type TextResult struct {
	TokenCount       int32  `json:"tokenCount"`
	OutputText       string `json:"outputText"`
	CompletionReason string `json:"completionReason"`
}

type TextResponse struct {
	InputTextTokenCount int32        `json:"inputTextTokenCount"`
	Results             []TextResult `json:"results"`
}

type TextGenerationConfig struct {
	// Use a lower value to decrease randomness in responses.
	Temperature float64 `json:"temperature"`
	// Use a lower value to ignore less probable options and decrease the diversity of responses.
	TopP float64 `json:"topP"`
	// Specify the maximum number of tokens to generate in the response. Maximum token limits are strictly enforced.
	MaxTokenCount int32 `json:"maxTokenCount"`
	// Specify a character sequence to indicate where the model should stop.
	StopSequences []string `json:"stopSequences,omitempty"`
}

type TextRequest struct {
	// The prompt to provide the model for generating a response.
	// To generate responses in a conversational style, wrap the prompt by using the following format:
	//
	// "inputText": "User: <prompt>\nBot:
	InputText string `json:"inputText"`
	// Optional, used to configure inference parameters
	TextGenerationConfig *TextGenerationConfig `json:"textGenerationConfig"`
}

type InvocationMetrics struct {
	InputTokenCount   int `json:"inputTokenCount"`
	OutputTokenCount  int `json:"outputTokenCount"`
	InvocationLatency int `json:"invocationLatency"`
	FirstByteLatency  int `json:"firstByteLatency"`
}

type TextStreamedResponse struct {
	OutputText                string            `json:"outputText"`
	Index                     int               `json:"index"`
	TotalOutputTextTokenCount int               `json:"totalOutputTextTokenCount"`
	CompletionReason          string            `json:"completionReason"`
	InputTextTokenCount       int               `json:"inputTextTokenCount"`
	InvocationMetrics         InvocationMetrics `json:"amazon-bedrock-invocationMetrics"`
}
