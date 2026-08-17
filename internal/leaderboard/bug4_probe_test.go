package leaderboard

import (
	"runtime"
	"sync"
	"testing"
)

func TestProbe_ConcurrentTopNIsRaceFree(t *testing.T) {
	runtime.GOMAXPROCS(4)
	b := New()
	if _, _, _, err := b.Submit("alice", 0, 0); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 32; i++ {
		if _, _, _, err := b.Submit(string(rune('a'+i)), int64(i), int64(i)); err != nil {
			t.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		for i := 1; i <= 2000; i++ {
			b.Submit("alice", int64(i), int64(i))
		}
	}()
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 2000; j++ {
				b.TopN(10)
			}
		}()
	}
	wg.Wait()
}
