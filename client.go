package aimagician

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

	"golang.org/x/net/websocket"
)

type Config struct {
}

type Client struct {
	Cookies        []http.Cookie
	ConversationId string
	TaskId         int32
}

type ChatStreamResponse struct {
	Conn *websocket.Conn
}

func (c *ChatStreamResponse) Recv() (*ChatResponse, error) {
	buf := make([]byte, 1024)
	len, err := c.Conn.Read(buf)
	if err != nil {
		return nil, err
	}

	resMsg := ChatResponse{}
	err = json.Unmarshal(buf[:len], &resMsg)
	if err != nil {
		return nil, err
	}
	if resMsg.Action == "end" || resMsg.Action == "error" {
		return nil, io.EOF
	}
	return &resMsg, nil
}

func (c *ChatStreamResponse) Close() {
	c.Conn.Close()
}

const (
	ws_url           = "wss://ai-prompt-chat-apis.guokr.net/chat"
	tasks_url        = "https://ai-prompt.guokr.net/api/recommend/tasks"
	conversation_url = "https://ai-prompt.guokr.net/api/chat/conversation"
)

func NewClient(cookies []http.Cookie) *Client {
	return &Client{
		Cookies: cookies,
	}
}

func (c *Client) randomTaskId() int32 {
	request, err := http.NewRequest("GET", tasks_url, nil)
	if err != nil {
		return -1
	}
	for _, cookie := range c.Cookies {
		request.AddCookie(&cookie)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return -1
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1
	}
	tasks := TasksResponse{}
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return -1
	}
	if len(tasks.Tasks) == 0 {
		return -1
	}
	idx := rand.Int31n(int32(len(tasks.Tasks)))
	return tasks.Tasks[idx].TaskId
}

func (c *Client) newConversationId() string {
	request, err := http.NewRequest("GET", conversation_url, nil)
	if err != nil {
		return ""
	}
	for _, cookie := range c.Cookies {
		request.AddCookie(&cookie)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	conversation := ConverstaionResponse{}
	err = json.Unmarshal(data, &conversation)
	if err != nil {
		return ""
	}
	return conversation.ConversationId
}

func (c *Client) ChatStream(content string) *ChatStreamResponse {
	if c.TaskId == 0 {
		c.TaskId = c.randomTaskId()
	}
	if c.ConversationId == "" {
		c.ConversationId = c.newConversationId()
	}

	conn, err := websocket.Dial(ws_url, "", ws_url)
	if err != nil {
		return nil
	}
	msg := ChatRequest{
		AI_Persona_Id:  25,
		ConversationId: c.ConversationId,
		TaskId:         c.TaskId,
		Content:        content,
	}
	data, err := json.Marshal(&msg)
	if err != nil {
		return nil
	}
	_, err = conn.Write(data)
	if err != nil {
		return nil
	}
	return &ChatStreamResponse{
		Conn: conn,
	}
}

func (c *Client) Chat(content string) string {
	if c.TaskId == 0 {
		c.TaskId = c.randomTaskId()
	}
	if c.ConversationId == "" {
		c.ConversationId = c.newConversationId()
	}

	conn, err := websocket.Dial(ws_url, "", ws_url)
	defer conn.Close()
	if err != nil {
		return ""
	}
	msg := ChatRequest{
		AI_Persona_Id:  25,
		ConversationId: c.ConversationId,
		TaskId:         c.TaskId,
		Content:        content,
	}
	data, err := json.Marshal(&msg)
	if err != nil {
		return ""
	}
	_, err = conn.Write(data)
	if err != nil {
		return ""
	}
	var res string
	var need bool = false
	for {
		buf := make([]byte, 1024)
		len, err := conn.Read(buf)
		if err != nil {
			break
		}

		resMsg := ChatResponse{}
		err = json.Unmarshal(buf[:len], &resMsg)
		if err != nil {
			break
		}
		if resMsg.Action != "normal" {
			continue
		}
		if resMsg.Residual == 0 && !need {
			need = true
		}

		res += resMsg.Content
		if resMsg.Action == "end" {
			break
		}
	}
	if need {
		c.ConversationId = c.newConversationId()
	}
	return res
}
