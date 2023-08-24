package main

import (
	"fmt"
	"sync"

	"golang.org/x/tour/tree"
)

type counter struct {
	mu sync.Mutex
	c  int
}

func Walk(t *tree.Tree, ch, quit chan int, c *counter, wg *sync.WaitGroup) {
	if t != nil {
		c.mu.Lock()
		fmt.Printf("out: c.c: %v & ", c.c)
		if c.c == 0 {
			c.mu.Unlock()
			fmt.Printf("in: c.c: %v\n", c.c)
			quit <- 1
			close(ch)
			wg.Done()
			return
		}
		c.c--

		c.mu.Unlock()

		select {
		case <-quit:
			fmt.Printf("quit: i got it")
			return
		case ch <- t.Value:
			fmt.Printf("t.Value: %v\n", t.Value)
			wg.Add(2)
			go Walk(t.Left, ch, quit, c, wg)
			go Walk(t.Right, ch, quit, c, wg)
		}
	}

	wg.Done()
}

func Same(t1 *tree.Tree) {
	ch := make(chan int, 10)
	quit := make(chan int)
	c := &counter{c: 10}
	var wg sync.WaitGroup
	wg.Add(1)
	go Walk(t1, ch, quit, c, &wg)
	go func() {
		defer wg.Done()
		j := 0
		for i := range ch {
			fmt.Printf("%v, i: %v\n", j, i)
			j++
		}
	}()
	wg.Wait()

}

func main() {
	t0 := tree.New(1)

	fmt.Println(t0)
	Same(t0)

}
