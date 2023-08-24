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
		defer c.mu.Unlock()
		defer wg.Done()
		fmt.Printf("out: c.c: %v & ", c.c)
		if c.c--; c.c == 0 {
			fmt.Printf("in: c.c: %v\n", c.c)
			quit <- 1
			close(ch)
			return
		}

		select {
		case <-quit:
			fmt.Printf("quit: i got it")
			return
		case ch <- t.Value:
			fmt.Printf("t.Value: %v\n", t.Value)
			go Walk(t.Left, ch, quit, c, wg)
			go Walk(t.Right, ch, quit, c, wg)
		}
	}

}

func Same(t1 *tree.Tree) {
	ch := make(chan int, 10)
	quit := make(chan int)
	c := &counter{c: 10}
	var wg sync.WaitGroup
	wg.Add(10)
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
