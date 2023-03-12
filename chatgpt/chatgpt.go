package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/linzb0123/chatgpt-demo/utils"
	"io"
	"log"
	"net/http"
)

const (
	DefaultModel  = "gpt-3.5-turbo-0301"
	OpenAiChatURL = "https://api.openai.com/v1/chat/completions"
)

type ChatGPT interface {
	Chat(content string) (string, error)
	Reset()
	GetModel() string
}

type chatGPT struct {
	apiKey   string
	model    string
	prompts  []*promptMessage
	question *promptMessage
	reply    *promptMessage
}

type promptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func New(apiKey string, model string) ChatGPT {
	return &chatGPT{
		apiKey: apiKey,
		model:  utils.If(model == "", DefaultModel, model),
	}
}

func (c *chatGPT) Reset() {
	c.prompts = nil
	c.question = nil
	c.reply = nil
}

func (c *chatGPT) GetModel() string {
	return c.model
}

func (c *chatGPT) Chat(content string) (string, error) {

	if len(content) == 0 {
		return "", errors.New("content empty")
	}

	c.question = &promptMessage{
		Role:    "user",
		Content: content,
	}

	resp, err := c.request()
	if err != nil {
		log.Printf("chat request err:%v", err)
		return "", err
	}

	if len(resp.Choices) == 0 {
		log.Printf("resp not chices.resp body:%s", utils.Dump(resp))
		return "请稍后重试", nil
	}

	c.reply = &resp.Choices[0].Message

	c.prompts = append(c.prompts, c.question)
	c.prompts = append(c.prompts, c.reply)

	return c.reply.Content, nil
}

type requestBody struct {
	Model    string           `json:"model"`
	Messages []*promptMessage `json:"messages"`
}
type respBody struct {
	Id      string            `json:"id,omitempty"`
	Object  string            `json:"object,omitempty"`
	Created int               `json:"created,omitempty"`
	Model   string            `json:"model,omitempty"`
	Usage   respBodyUsage     `json:"usage"`
	Choices []respBodyChoices `json:"choices,omitempty"`
}

type respBodyUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type respBodyChoices struct {
	Message      promptMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
	Index        int           `json:"index"`
}

func (c *chatGPT) request() (*respBody, error) {
	prompts := append(c.prompts, c.question)

	body := requestBody{
		Model:    c.model,
		Messages: prompts,
	}

	reqBody := bytes.NewBuffer(utils.Dump(body))
	req, err := http.NewRequest(http.MethodPost, OpenAiChatURL, reqBody)
	if err != nil {
		log.Printf("new request body[%s] err:%v", reqBody.Bytes(), err)
		return nil, err
	}

	//req.Close = true

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	//log.Printf("raw req body:%s", reqBody.Bytes())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("resp err:%v", err)
		return nil, err
	}

	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("io read err:%v", err)
		return nil, err
	}

	log.Printf("raw resp body:%s", bs)

	var rep respBody
	err = json.Unmarshal(bs, &rep)
	if err != nil {
		log.Printf("unmarshal err:%v body:%s", err, bs)
		return nil, err
	}

	return &rep, nil
}
