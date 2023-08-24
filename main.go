package main

import (
	"fmt"
	"sort"
	"sync"

	"golang.org/x/exp/slices"
	"golang.org/x/tour/tree"
)

type counter struct {
	mu sync.Mutex
	c  int
}

func Walk(t *tree.Tree, ch, quit chan int, c *counter, wg *sync.WaitGroup) {
	if t != nil {
		c.mu.Lock()
		if c.c == 0 {
			c.mu.Unlock()
			quit <- 1
			close(ch)
			wg.Done()
			return
		}
		c.c--

		c.mu.Unlock()

		select {
		case <-quit:
			return
		case ch <- t.Value:
			wg.Add(2)
			go Walk(t.Left, ch, quit, c, wg)
			go Walk(t.Right, ch, quit, c, wg)
		}
	}

	wg.Done()
}

func array(t *tree.Tree, a *[]int, w *sync.WaitGroup) {
	ch := make(chan int, 10)
	quit := make(chan int)
	c := &counter{c: 10}
	var wg sync.WaitGroup
	wg.Add(1)
	go Walk(t, ch, quit, c, &wg)
	go func() {
		defer wg.Done()
		for i := range ch {
			*a = append(*a, i)

		}
	}()
	wg.Wait()
	w.Done()

}
func Same(t1, t2 *tree.Tree) bool {
	var a, b []int
	var wg sync.WaitGroup
	wg.Add(2)
	go array(t1, &a, &wg)
	go array(t2, &b, &wg)
	wg.Wait()
	sort.Ints(a)
	sort.Ints(b)

	return slices.Equal(a, b)
}

func main() {
	t1 := tree.New(1)
	t2 := tree.New(1)

	fmt.Printf("Same(t1, t2): %v\n", Same(t1, t2))
	t3 := tree.New(1)
	t4 := tree.New(2)

	fmt.Printf("Same(t1, t2): %v\n", Same(t3, t4))

}
