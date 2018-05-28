package crawler

import "sync"

type WordMap struct {
	mu       sync.RWMutex
	internal map[string]int
}

func (w *WordMap) Add(word string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.internal[word]; ok {
		w.internal[word]++
	} else {
		w.internal[word] = 1
	}

}

func (w *WordMap) Sort() {
	// Implement clever solution
	w.mu.Lock()
	defer w.mu.Unlock()

	return
}

func NewWordMap() *WordMap {
	return &WordMap{
		internal: make(map[string]int),
	}
}
