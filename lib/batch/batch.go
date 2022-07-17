package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) []user {
	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		res    = make([]user, 0, n)
		reduce = make(chan struct{}, pool)
	)

	var i int64
	for i = 0; i < n; i++ {
		reduce <- struct{}{}
		wg.Add(1)
		go func(i int64) {
			user := getOne(i)
			mu.Lock()
			res = append(res, user)
			mu.Unlock()
			<-reduce
			wg.Done()
		}(i)
	}

	wg.Wait()
	return res
}
