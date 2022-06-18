package quiz

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type QA struct{ Q, A string }

type Quizzer struct {
	booklet []QA
	mu      sync.RWMutex
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewQuizzer(booklet []QA) *Quizzer {
	var q Quizzer
	q.booklet = booklet
	return &q
}

func (q *Quizzer) Length() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.booklet)
}

func (q *Quizzer) AppendQA(qa QA) {
	q.mu.Lock()
	q.booklet = append(q.booklet, qa)
	q.mu.Unlock()
}

func (q *Quizzer) Conduct(d time.Duration, shuf bool) (score int) {
	scanner := bufio.NewScanner(os.Stdin)

	l := q.Length()
	ord := make([]int, l)
	for i := 0; i < l; i++ {
		ord[i] = i
	}
	if shuf {
		rand.Shuffle(l, func(i, j int) {
			ord[i], ord[j] = ord[j], ord[i]
		})
	}

	done := make(chan bool)
	go func() {
		for i, qi := range ord {
			q.mu.RLock()
			fmt.Printf("Problem #%v: %v = ", i+1, q.booklet[qi].Q)
			q.mu.RUnlock()

			if b := scanner.Scan(); !b {
				fmt.Println()
				break
			}
			res := scanner.Text()

			if q.equalAnswer(qi, res) {
				score++
			}
		}
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(d):
		fmt.Println()
	}
	return score
}

func (q *Quizzer) equalAnswer(i int, res string) bool {
	q.mu.RLock()
	ansf := strings.Fields(q.booklet[i].A)
	q.mu.RUnlock()
	resf := strings.Fields(res)

	return equalFields(ansf, resf)
}

func equalFields(resf, ansf []string) bool {
	if len(resf) != len(ansf) {
		return false
	}

	for i := range resf {
		if !strings.EqualFold(resf[i], ansf[i]) {
			return false
		}
	}
	return true
}
