package aimagician

type VerificationRequest struct {
	Phone       string `josn:"phone"`
	Action_Type string `json:"action_type" default:"signin"`
}

type VerificationResponse struct {
	Status    int    `json:"status"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	errors    string
}

type LoginRequest struct {
	Phone       string `json:"phone"`
	Verify_Code string `json:"verify_code"`
}

type TaskInfo struct {
	Count  int32  `json:"count"`
	TaskId int32  `json:"task_id"`
	Title  string `json:"title"`
}

type TasksResponse struct {
	Scene string     `json:"scene"`
	Tasks []TaskInfo `json:"tasks"`
}

type ConverstaionResponse struct {
	MaxChatCount   int32  `json:"max_chat_count"`
	MaxTokens      int32  `json:"max_tokens"`
	ConversationId string `json:"conversation_id"`
}

type ChatRequest struct {
	TaskId         int32  `json:"task_id"`
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content"`
}

type ChatResponse struct {
	Task_id         int32  `json:"task_id"`
	Tonversation_id string `json:"conversation_id"`
	Chat_id         string `json:"chat_id"`
	Content         string `json:"content"`
	Action          string `json:"action"`
	MsgType         string `json:"type"`
	Residual        int32  `json:"residual"`
}
