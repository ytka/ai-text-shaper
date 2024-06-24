package openai

type TotalUsageCost struct {
	UsageCosts []*UsageCost
}

func NewTotalUsageCost(usageCosts []*UsageCost) *TotalUsageCost {
	return &TotalUsageCost{
		UsageCosts: usageCosts,
	}
}

func (tuc *TotalUsageCost) TotalPromptTokens() int {
	totalPromptTokens := 0
	for _, uc := range tuc.UsageCosts {
		totalPromptTokens += uc.PromptTokens()
	}
	return totalPromptTokens
}

func (tuc *TotalUsageCost) TotalCompletionTokens() int {
	totalCompletionTokens := 0
	for _, uc := range tuc.UsageCosts {
		totalCompletionTokens += uc.CompletionTokens()
	}
	return totalCompletionTokens
}

func (tuc *TotalUsageCost) TotalTotalTokens() int {
	totalTokens := 0
	for _, uc := range tuc.UsageCosts {
		totalTokens += uc.TotalTokens()
	}
	return totalTokens
}

func (tuc *TotalUsageCost) TotalPromptTokensCost() (bool, float64) {
	var totalCost, cost float64
	var ok bool
	for _, uc := range tuc.UsageCosts {
		if ok, cost = uc.PromptTokensCost(); !ok {
			return false, 0
		}
		totalCost += cost
	}
	return true, totalCost
}

func (tuc *TotalUsageCost) TotalCompletionTokensCost() (bool, float64) {
	var totalCost, cost float64
	var ok bool
	for _, uc := range tuc.UsageCosts {
		if ok, cost = uc.CompletionTokensCost(); !ok {
			return false, 0
		}
		totalCost += cost
	}
	return true, totalCost
}

func (tuc *TotalUsageCost) TotalTotalTokensCost() (bool, float64) {
	var totalCost, cost float64
	var ok bool
	for _, uc := range tuc.UsageCosts {
		if ok, cost = uc.TotalTokensCost(); !ok {
			return false, 0
		}
		totalCost += cost
	}
	return true, totalCost
}
