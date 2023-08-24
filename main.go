package main

import (
	"fmt"
	"sort"
	"sync"

	"golang.org/x/exp/slices"
	"golang.org/x/tour/tree"
)

func Walk(t *tree.Tree, ch chan int, wg *sync.WaitGroup) {
	if t != nil {
		ch <- t.Value
		wg.Add(2)
		go Walk(t.Left, ch, wg)
		go Walk(t.Right, ch, wg)

	}

	wg.Done()
}

func array(t *tree.Tree, a *[]int, w *sync.WaitGroup) {
	ch := make(chan int, 10)
	var wg sync.WaitGroup
	wg.Add(1)
	go Walk(t, ch, &wg)
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
	t3 := tree.New(2)

	fmt.Printf("Same(t1, t2): %v\n", Same(t2, t3))

}
