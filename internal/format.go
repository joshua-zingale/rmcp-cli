package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/joshua-zingale/remote-mcp-host/remote-mcp-host/api"
)

func FormatServerList(list *api.McpServerList) string {
	bldr := strings.Builder{}

	bldr.WriteString("## Server List ##\n")
	sort.Slice(list.Servers, func(i, j int) bool {
		return list.Servers[i].Name < list.Servers[j].Name
	})
	for i, server := range list.Servers {
		bldr.WriteString(fmt.Sprintf("  %d. %s\n", i+1, server.Name))
	}

	return bldr.String()
}

func FormatGenerationResponse(resp *api.GenerationResponse) string {
	bldr := strings.Builder{}

	bldr.WriteString("## Generation Response ##\n")

	bldr.WriteString(fmt.Sprintf("Role: %s\n", resp.Message.Role))

	for i, part := range resp.Message.Parts {
		bldr.WriteString(fmt.Sprintf("Part %d: %s", i+1, FormatPart(&part)))
	}

	return bldr.String()
}

func FormatPart(part *api.UnionPart) string {
	bldr := strings.Builder{}
	bldr.WriteString("Part: ")
	switch part := part.Part.(type) {
	case api.TextPart:
		bldr.WriteString("Text: ")
		if part.Error != "" {
			bldr.WriteString("Error!: ")
			bldr.WriteString(part.Error)
		} else {
			bldr.WriteString(part.Text)
		}
	case api.ToolUsePart:
		bldr.WriteString("ToolUse: ")
		bldr.WriteString(fmt.Sprintf("ToolID: %+v : ", part.ToolId))
		bldr.WriteString(fmt.Sprintf("Input: %+v : ", part.Input))

		if part.Error != "" {
			bldr.WriteString(fmt.Sprintf("Error!: %s", part.Error))
		} else {
			bldr.WriteString(fmt.Sprintf("Output: %+v", part.Output))
		}
	default:
		bldr.WriteString("Error!: Unknown part type encountered.")
	}

	return bldr.String()

}
