package openai

type TotalUsageCost struct {
	UsageCosts []*UsageCost
}

// NewTotalUsageCost creates a new TotalUsageCost instance.
func NewTotalUsageCost(usageCosts []*UsageCost) *TotalUsageCost {
	return &TotalUsageCost{
		UsageCosts: usageCosts,
	}
}

// TotalPromptTokens calculates the total number of prompt tokens.
func (tuc *TotalUsageCost) TotalPromptTokens() int {
	totalPromptTokens := 0
	for _, uc := range tuc.UsageCosts {
		totalPromptTokens += uc.PromptTokens()
	}
	return totalPromptTokens
}

// TotalCompletionTokens calculates the total number of completion tokens.
func (tuc *TotalUsageCost) TotalCompletionTokens() int {
	totalCompletionTokens := 0
	for _, uc := range tuc.UsageCosts {
		totalCompletionTokens += uc.CompletionTokens()
	}
	return totalCompletionTokens
}

// TotalTotalTokens calculates the total number of tokens.
func (tuc *TotalUsageCost) TotalTotalTokens() int {
	totalTokens := 0
	for _, uc := range tuc.UsageCosts {
		totalTokens += uc.TotalTokens()
	}
	return totalTokens
}

// TotalPromptTokensCost calculates the total cost of prompt tokens.
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

// TotalCompletionTokensCost calculates the total cost of completion tokens.
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

// TotalTotalTokensCost calculates the total cost of tokens.
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
