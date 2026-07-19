package bedrock

func NewRequestBody(promptText string, b64Data string) RequestBody {
	return RequestBody{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        1024,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{Type: "text", Text: promptText},
					{
						Type: "image",
						Source: &ImageSource{
								Type:      "base64",
								MediaType: "image/jpeg",
								Data:      b64Data,
						},
					},
				},
			},
		},
	}
}
