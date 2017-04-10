package u

// Semaphore based on http://golang.org/doc/effective_go.html#channels
type Semaphore struct {
	sem chan bool
}

// NewSemaphore creates a semaphore
func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		sem: make(chan bool, max),
	}
}

// Enter enters a semaphore
func (s *Semaphore) Enter() {
	s.sem <- true
}

// Leave leaves a semaphore
func (s *Semaphore) Leave() {
	<-s.sem
}
