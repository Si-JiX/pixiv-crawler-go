package threadpool

import (
	"fmt"
	"pixiv-cil/pkg/config"
	"sync"
)

var Threading *ThreadStruct

type ThreadStruct struct {
	wg             *sync.WaitGroup
	ch             chan struct{}
	ProgressLength int
	progressCount  int
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
func (t *ThreadStruct) Len() int {
	return len(t.ch)
}

func (t *ThreadStruct) ProgressCountAdd() {
	t.progressCount += 1
}

func (t *ThreadStruct) GetProgressInfo() {
	fmt.Printf("download image:%d/%d\r", Threading.progressCount, Threading.ProgressLength)
}

func (t *ThreadStruct) Close() {
	t.Wait()
	t.progressCount = 0
	t.ProgressLength = 0
}
