package internal

import (
	"context"
	"fmt"

	"github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/event"
	opt "go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo = new(mongo)

type mongo struct{}

func (m *mongo) GetClientOptions() []options.ClientOptions {
	cmdMonitor := &event.CommandMonitor{
		Started: func(ctx context.Context, event *event.CommandStartedEvent) {
			logx.Info(fmt.Sprintf("[MongoDB][RequestID:%d][database:%s] %s\n", event.RequestID, event.DatabaseName, event.Command))
		},
		Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
			logx.Info(fmt.Sprintf("[MongoDB][RequestID:%d] [%s] %s\n", event.RequestID, event.Duration.String(), event.Reply))
		},
		Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
			logx.Error(fmt.Sprintf("[MongoDB][RequestID:%d] [%s] %s\n", event.RequestID, event.Duration.String(), event.Failure))
		},
	}
	return []options.ClientOptions{{ClientOptions: &opt.ClientOptions{Monitor: cmdMonitor}}}
}
