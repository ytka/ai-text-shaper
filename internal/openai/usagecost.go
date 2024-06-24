package openai

type UsageCost struct {
	chatCompletion *ChatCompletion
}

func NewUsageCost(chatCompletion *ChatCompletion) *UsageCost {
	return &UsageCost{chatCompletion: chatCompletion}
}

func (uc *UsageCost) ModelName() string {
	return uc.chatCompletion.Model
}

func (uc *UsageCost) PromptTokens() int {
	return uc.chatCompletion.Usage.PromptTokens
}

func (uc *UsageCost) CompletionTokens() int {
	return uc.chatCompletion.Usage.CompletionTokens
}

func (uc *UsageCost) TotalTokens() int {
	return uc.chatCompletion.Usage.TotalTokens
}

func (uc *UsageCost) PromptTokensCost() (bool, float64) {
	return CalculateInputTokensCost(uc.ModelName(), float64(uc.PromptTokens()))
}

func (uc *UsageCost) CompletionTokensCost() (bool, float64) {
	return CalculateOutputTokensCost(uc.ModelName(), float64(uc.PromptTokens()))
}

func (uc *UsageCost) TotalTokensCost() (bool, float64) {
	var totalCost, cost float64
	var ok bool

	if ok, cost = uc.PromptTokensCost(); !ok {
		return false, 0
	}
	totalCost += cost

	if ok, cost = uc.CompletionTokensCost(); !ok {
		return false, 0
	}
	totalCost += cost

	return true, totalCost
}
