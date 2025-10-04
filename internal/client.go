package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (client *Client) Generate(req *rmcp.GenerationRequest, opts *GenerateOpts) (*rmcp.GenerationResponse, error) {

	resp, err := post[rmcp.GenerationResponse](client, "/generations", req)
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

// Sends a post request to the MCP host. The bodyObj should be nil if there is no body
func post[T any](client *Client, path string, bodyObj any) (*T, error) {

	var reqBody io.Reader

	if bodyObj == nil {
		reqBody = nil
	} else {
		bytesBody, err := json.Marshal(bodyObj)
		if err != nil {
			return nil, fmt.Errorf("could not marshal object: %s", err)
		}
		reqBody = bytes.NewReader(bytesBody)
	}

	req, err := http.NewRequest("POST", client.address+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("problem forming request: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("problem sending request: %s", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("invalid request (%d): %s", resp.StatusCode, body)
	}

	var responseObj T
	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		return nil, fmt.Errorf("problem decoding response: %s", err)
	}

	return &responseObj, err
}
