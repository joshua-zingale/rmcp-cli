package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	rmcp "github.com/joshua-zingale/remote-mcp-host/remote-mcp-host"
)

type Client struct {
	address string
	http    http.Client
}

func NewClient(host string, port string) Client {
	return Client{
		address: host + ":" + port,
		http:    http.Client{},
	}
}

type UserMessageOpts struct{}

func NewUserMessage(body string, opts *UserMessageOpts) rmcp.Message {
	return rmcp.Message{
		Role:  "user",
		Parts: []rmcp.UnionPart{{Part: rmcp.NewTextPart(body)}},
	}
}

func (client *Client) ListServers() (*rmcp.McpServerList, error) {

	list, err := get[rmcp.McpServerList](client, "/servers")
	if err != nil {
		return nil, err
	}
	return list, nil
}

type GenerateOpts struct{}

func (client *Client) Generate(messages []rmcp.Message, opts *GenerateOpts) (*rmcp.GenerationResponse, error) {

	req := rmcp.GenerationRequest{
		Messages: messages,
	}
	resp, err := post[rmcp.GenerationResponse](client, "/generations", &req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Sends a get request to the MCP host
func get[T any](client *Client, path string) (*T, error) {

	req, err := http.NewRequest("GET", client.address+"/"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, err
	}

	var responseObj T
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		return nil, err
	}

	return &responseObj, err
}

// Sends a post request to the MCP host. The bodyObj should be null if there is no body
func post[T any](client *Client, path string, bodyObj any) (*T, error) {

	var reqBody io.Reader

	if bodyObj == nil {
		reqBody = nil
	} else {
		bytesBody, err := json.Marshal(bodyObj)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(bytesBody)
	}

	req, err := http.NewRequest("POST", client.address+"/"+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, err
	}

	var responseObj T
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		return nil, err
	}

	return &responseObj, err
}
