package wait

import "sync"

func ForWaitGroupDone(wg *sync.WaitGroup) <-chan struct{} {
	done := make(chan struct{})

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(done)
	}(wg)

	return done
}
