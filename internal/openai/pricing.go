package openai

type Pricing struct {
	Model                  string
	InputTokensCostDollar  float64
	InputTokens            float64
	OutputTokensCostDollar float64
	OutputTokens           float64
}

const OneMillion = 1_000_000

var pricingList = []Pricing{
	{Model: "gpt-4o", InputTokensCostDollar: 5, InputTokens: OneMillion, OutputTokensCostDollar: 15, OutputTokens: OneMillion},
	{Model: "gpt-4o-2024-05-13", InputTokensCostDollar: 5, InputTokens: OneMillion, OutputTokensCostDollar: 15, OutputTokens: OneMillion},
	{Model: "gpt-3.5-turbo-0125", InputTokensCostDollar: 0.5, InputTokens: OneMillion, OutputTokensCostDollar: 1.5, OutputTokens: OneMillion},
	{Model: "gpt-3.5-turbo-instruct", InputTokensCostDollar: 1.5, InputTokens: OneMillion, OutputTokensCostDollar: 2, OutputTokens: OneMillion},
}

var pricingMap = map[string]Pricing{}

func init() {
	for _, p := range pricingList {
		pricingMap[p.Model] = p
	}
}

func GetPricing(model string) *Pricing {
	if p, ok := pricingMap[model]; ok {
		return &p
	}
	return nil
}

func CalculateInputTokensCost(model string, inputTokens float64) (bool, float64) {
	p := GetPricing(model)
	if p == nil {
		return false, 0
	}
	return true, p.InputTokensCostDollar * inputTokens / p.InputTokens
}

func CalculateOutputTokensCost(model string, outputTokens float64) (bool, float64) {
	p := GetPricing(model)
	if p == nil {
		return false, 0
	}
	return true, p.OutputTokensCostDollar * outputTokens / p.OutputTokens
}
