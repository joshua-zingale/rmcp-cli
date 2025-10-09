package internal

import (
	"regexp"

	"github.com/joshua-zingale/remote-mcp-host/remote-mcp-host/api"
)

var toolConfigs = make(map[api.ToolId]map[string]any)

var argRegex = regexp.MustCompile(`[^\s]+`)

func HandleCommand(co CommandOutput, ci CommandInput) {
	cmd := ci.matchGroups[1]
	args := argRegex.FindAllString(ci.matchGroups[2], -1)

	switch cmd {
	case "servers":
		r, err := ci.client.ListServers()
		if err != nil {
			co.Printf("%s\n", err)
			return
		}
		co.Print(FormatServerList(r))
	case "patch":
		if len(args) != 4 {
			co.Printf("/patch does not take %d, but 4 arguments: must be: /path SERVERNAME TOOLNAME PARAMETER ARGUMENT\n", len(args))
			return
		}
		toolId := api.ToolId{ServerName: args[0], Name: args[1]}
		parameter := args[2]
		argument := args[3]

		if hash, ok := toolConfigs[toolId]; ok {
			hash[parameter] = argument
		} else {
			toolConfigs[toolId] = map[string]any{parameter: argument}
		}
		co.Printf("Added patch to %+v: %s -> %s\n", toolId, parameter, argument)
	case "list":
		if len(args) != 1 {
			co.Printf("/list takes 1 argument\n")
			return
		}
		switch args[0] {
		case "patch":
			co.Printf("%+v\n", toolConfigs)
		}
	default:
		co.Printf("Invalid command: '%s'\n", cmd)

	}
}
func HandleUserMessage(co CommandOutput, ci CommandInput) {
	message := ci.Match()

	var tcfgs []api.ToolConfig
	for id, input := range toolConfigs {
		tcfgs = append(tcfgs, api.ToolConfig{
			ToolId: id,
			ToolPatch: api.ToolPatch{
				Input: input,
			},
		})
	}
	req := api.GenerationRequest{
		Messages:    []api.Message{NewUserMessage(message, nil)},
		ToolConfigs: tcfgs,
	}

	res, err := ci.client.Generate(&req, nil)
	if err != nil {
		co.Printf("ERROR: %s\n", err)
		return
	}
	co.Printf("%s\n", FormatGenerationResponse(res))
}
