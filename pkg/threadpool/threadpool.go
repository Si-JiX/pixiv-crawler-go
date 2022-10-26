package threadpool

import (
	"pixiv-cil/pkg/config"
	"sync"
)

var Threading *ThreadStruct

type ThreadStruct struct {
	wg *sync.WaitGroup
	ch chan struct{}
}

func InitThread() {
	if config.Vars.ThreadMax == 0 {
		config.Vars.ThreadMax = 16
	}
	Threading = &ThreadStruct{wg: &sync.WaitGroup{}, ch: make(chan struct{}, config.Vars.ThreadMax)}
}

func (t *ThreadStruct) Add() {
	t.wg.Add(1)
	t.ch <- struct{}{}
}
func (t *ThreadStruct) Done() {
	t.wg.Done()
	<-t.ch
}

func (t *ThreadStruct) Wait() {
	t.wg.Wait()

}
func (t *ThreadStruct) Close() {
	close(t.ch)
}
func (t *ThreadStruct) Len() int {
	return len(t.ch)
}
