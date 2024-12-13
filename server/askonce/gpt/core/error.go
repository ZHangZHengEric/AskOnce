package core

type OpenAIError struct {
	Code    any    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Param   string `json:"param,omitempty"`
}

func (err OpenAIError) Error() string {
	return err.Message
}

func NewOpenAiError(code any, message string) OpenAIError {
	return OpenAIError{
		Code:    code,
		Message: message,
	}
}
