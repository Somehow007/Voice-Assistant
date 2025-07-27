package tools

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

var AllTools = []Tool{
	{
		Type: "function",
		Function: Function{
			Name:        "get_current_time",
			Description: "当用户查询现在或者说当前时间时，请调用这个方法。",
			Parameters:  map[string]interface{}{},
		},
	},
	{
		Type: "function",
		Function: Function{
			Name:        "get_current_weather",
			Description: "当用户询问指定城市的天气时请使用这个方法。",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "城市或县区，比如北京市、杭州市、余杭区等。",
					},
				},
				"required": []string{"location"},
			},
		},
	},
}
