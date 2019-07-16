package k8s_test

import (
	"testing"
)

func TestChan(t *testing.T) {

	ch := make(chan int)

	go func() {

		defer func() {
			if r := recover(); r != nil {
				t.Log(r)
			}
		}()

		for i:=0; i<10; i++ {
			ch <- i
		}
		close(ch)
		ch <- 101
	}()

	for {
		if num, ok := <-ch; ok {
			t.Log(num)
			continue
		}
		t.Log("end")
		break
	}
}