package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	eh "github.com/looplab/eventhorizon"
)

func init() {
	eh.RegisterCommand(func() eh.Command { return &GET{} })
}

const (
	GETCommand = eh.CommandType("http:RedfishResource:GET")
)

// Static type checking for commands to prevent runtime errors due to typos
var _ = eh.Command(&GET{})

// HTTP GET Command
type GET struct {
	ID      eh.UUID           `json:"id"`
	CmdID   eh.UUID           `json:"cmdid"`
	Headers map[string]string `eh:"optional"`
}

func (c *GET) AggregateType() eh.AggregateType { return AggregateType }
func (c *GET) AggregateID() eh.UUID            { return c.ID }
func (c *GET) CommandType() eh.CommandType     { return GETCommand }
func (c *GET) SetAggID(id eh.UUID)             { c.ID = id }
func (c *GET) SetCmdID(id eh.UUID)             { c.CmdID = id }
func (c *GET) SetUserDetails(u string, p []string) string {
	return "checkMaster"
}
func (c *GET) Handle(ctx context.Context, a *RedfishResourceAggregate) error {
	fmt.Printf("HANDLE GET!\n")
	data := HTTPCmdProcessedData{
		CommandID:  c.CmdID,
		Results:    map[string]interface{}{},
		StatusCode: 200,
		Headers:    a.Headers,
	}

	a.ProcessMeta(ctx)

	for k, v := range a.Properties {
		if !strings.HasSuffix(k, "@meta") {
			data.Results[k] = v
		}
	}
	a.eventBus.HandleEvent(ctx, eh.NewEvent(HTTPCmdProcessed, data, time.Now()))
	return nil
}
