package chat

import (
	"sync"

	bpb "github.com/Bit-Nation/protobuffers"
)

// handles a set of protobuf messages
func (c *Chat) messagesHandler(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error) {

	wg := sync.WaitGroup{}
	if len(req.Messages) > 0 {
		for _, msg := range req.Messages {
			wg.Add(1)
			go func(msg *bpb.ChatMessage) {
				defer wg.Done()
				err := c.handleReceivedMessage(msg)
				if err != nil {
					logger.Error(err)
				}
			}(msg)
		}
		wg.Wait()
		return &bpb.BackendMessage_Response{}, nil
	}

	return nil, nil

}
