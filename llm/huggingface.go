package llm

type HuggingFace struct{
	Model string
	ApiKey string
}

func (h *HuggingFace) Name() string{
	return "Provider: HuggingFace Model: " +h.Model 
} 

func (h *HuggingFace) Generate(prompt string) (string,error){
	return "",nil
}
