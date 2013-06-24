package happening

import (
    "sync"
)

// An uninteresting service.
type Service struct {
    ch          chan bool
    waitGroup   *sync.WaitGroup
}


// Make a new Service.
func NewService() *Service {
    s := &Service{
        ch: make(chan bool),
        waitGroup: &sync.WaitGroup{},
    }
    s.waitGroup.Add(1)
    return s
}


// Stop the service by closing the service's channel.
// Blocks until the service is really stopped.
func (s *Service) Stop() {
    close(s.ch)
    s.waitGroup.Wait()
}