package internal

import rmcp "github.com/joshua-zingale/remote-mcp-host/remote-mcp-host"

func HandleCommand(co CommandOutput, ci CommandInput) {
	cmd := ci.Match()[1:]
	switch cmd {
	case "servers":
		r, err := ci.client.ListServers()
		if err != nil {
			co.Printf("%s\n", err)
			return
		}
		co.Print(FormatServerList(r))
	default:
		co.Printf("Invalid command: '%s'", cmd)

	}
}
func HandleUserMessage(co CommandOutput, ci CommandInput) {
	message := ci.Match()

	req := rmcp.GenerationRequest{
		Messages: []rmcp.Message{NewUserMessage(message, nil)},
	}

	res, err := ci.client.Generate(&req, nil)
	if err != nil {
		co.Printf("ERROR: %s\n", err)
		return
	}
	co.Printf("%s\n", FormatGenerationResponse(res))
}
