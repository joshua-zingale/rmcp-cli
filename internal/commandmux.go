package internal

import (
	"context"
	"fmt"
	"io"
	"regexp"
)

type CommandHandler interface {
	Handle(CommandOutput, CommandInput)
}

type CommandInput struct {
	ctx         context.Context
	matchGroups []string
	client      *Client
}

// the full text that was matched to arrive at the handler
func (ci *CommandInput) Match() string {
	return ci.matchGroups[0]
}

type CommandOutput struct {
	writer io.Writer
}

func (co *CommandOutput) Print(str string) error {
	_, err := co.writer.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}

func (co *CommandOutput) Printf(str string, a ...any) error {
	_, err := co.writer.Write([]byte(fmt.Sprintf(str, a...)))
	if err != nil {
		return err
	}
	return nil
}

type CommandMux struct {
	patterns []*regexp.Regexp
	handlers []CommandHandler
	writer   io.Writer
	client   *Client
}

func NewCommandMux(writer io.Writer, client *Client) CommandMux {
	return CommandMux{writer: writer, client: client}
}

func (cm *CommandMux) Handle(ctx context.Context, input string) error {
	for i, p := range cm.patterns {
		matches := p.FindStringSubmatch(input)
		if matches == nil {
			continue
		}

		cm.handlers[i].Handle(CommandOutput{writer: cm.writer}, CommandInput{ctx: ctx, matchGroups: matches, client: cm.client})
		return nil

	}

	return fmt.Errorf("'%s' did not match any input patterns", input)
}

func (cm *CommandMux) AddHandler(pattern string, handler CommandHandler) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	cm.patterns = append(cm.patterns, regex)
	cm.handlers = append(cm.handlers, handler)

	return nil
}

func (cm *CommandMux) AddHandlerFunc(pattern string, handlerFunc func(CommandOutput, CommandInput)) error {
	handler := FuncHandler{fn: handlerFunc}
	return cm.AddHandler(pattern, handler)
}

type FuncHandler struct {
	fn func(CommandOutput, CommandInput)
}

func (fh FuncHandler) Handle(co CommandOutput, ci CommandInput) {
	fh.fn(co, ci)
}
