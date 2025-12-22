package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//request struct
type chatRequest struct {
	Model 	 string 		`json:"model"`
	Messages []Message 		`json:"messages"`
	Stream 	 bool			`json:"stream"`
}

// LLM Client
type LLMClient struct {
	BaseURL string
	APIKey  string
	Model 	string
}
//Message
type Message struct {
	Role 	string `json:"role"`
	Content string `json:"content"`
}

//streamChunk is a fragment
type streamChunk struct {
	Content string
}
//LLMStream represent a stream response
type LLMStream struct {
	scanner *bufio.Scanner //add a buffer zone
	resp  	*http.Response
}

// LLM API return body
type sseResponse struct {
	Choices []struct {
		Delta struct {
			Content string 	`json:"content"`
		} `json:"delta"`
	}`json:"choices"`
}
//New LLmClient creat new LLMClient
func NewLLMClient(baseURL, apiKey,model string)*LLMClient {
	return &LLMClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
	}
}
//use it recv next chunk
func (s *LLMStream) Recv() (*streamChunk,error) {
	//read line by line and parse json body
	//1.read next line 
	if !s.scanner.Scan() {
		//no more Data
		if err := s.scanner.Err();err != nil{
			return nil,err
		}
		return nil, io.EOF
	}
	line := s.scanner.Text()
	//2. check is SSE data line
	if !strings.HasPrefix(line,"data: ") {
		return s.Recv()
	}

	//3. extract json part
	jsonStr := strings.TrimPrefix(line,"data: ")

	//4. checkout is done ?
	if jsonStr == "[DONE]" {
		return nil, io.EOF
	}

	//5. Parse json
	var resp sseResponse
	if err := json.Unmarshal([]byte(jsonStr),&resp); err != nil {
		return nil, fmt.Errorf("parser json failure: %w",err)
	}

	//6.extract content
	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Delta.Content
		return &streamChunk{Content: content},nil
	}
	//TODO 
	return s.Recv()
}
//close stream
func (s *LLMStream) Close() error {
	if s.resp != nil {
		return s.resp.Body.Close()
	}
	return nil
}

//send stream request
func (c *LLMClient) Stream (ctx context.Context,messages []Message) (*LLMStream,error) {
	// TODO: 第一步 - 构建请求体
    reqBody := chatRequest{
		Model:		c.Model,
		Messages: 	messages,
		Stream: 	true,	//开启流式返回
	}
    // TODO: 第二步 - 转换成 JSON
    jsonData,err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("转化 JSON failure: %v",err)
	}
    // TODO: 第三步 - 创建 HTTP 请求
    url := c.BaseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST",url,bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w",err)
	}
	// TODO: 第四步 - 设置请求头
    req.Header.Set("Authorization","Bearer " + c.APIKey)
	req.Header.Set("Content-Type", "application/json")
    // TODO: 第五步 - 发送请求
    client := &http.Client{}
	resp,err := client.Do(req)
	if err != nil {
		return nil,fmt.Errorf("发送请求失败: %w",err)
	}
    // TODO: 第六步 - 检查响应状态
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s",resp.StatusCode,resp.Status)
	}
    // TODO: 第七步 - 创建 LLMStream 返回
	scanner := bufio.NewScanner(resp.Body)

	return &LLMStream{
		scanner: scanner,
		resp:  	 resp,
	},nil

}