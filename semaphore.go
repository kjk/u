package u

// based on http://golang.org/doc/effective_go.html#channels
type Semaphore struct {
	sem chan int
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		sem: make(chan int, max),
	}
}

func (s *Semaphore) Enter() {
	s.sem <- 1
}

func (s *Semaphore) Leave() {
	<-s.sem
}
