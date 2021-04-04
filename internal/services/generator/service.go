package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const maxWaitTimeSec = 3

type Service struct {
	label string
	logs  chan Message
}

type Message struct {
	Label   string
	Text    string
	EventAt time.Time
}

func New(label string) *Service {
	return &Service{
		label: label,
		logs:  make(chan Message, 1),
	}
}

func (svc *Service) GenerateLogs(ctx context.Context) error {
	for {
		secs := rand.Intn(maxWaitTimeSec)
		m := genMsg(time.Duration(secs) * time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case svc.logs <- Message{
			Label:   svc.label,
			Text:    m,
			EventAt: time.Now().UTC(),
		}:
			println(svc.label, m)
		}
	}
}

func (svc *Service) Logs() chan Message {
	return svc.logs
}

func genMsg(waitTime time.Duration) string {
	time.Sleep(waitTime)
	return fmt.Sprintf("worker's wait time: %s", waitTime)
}
