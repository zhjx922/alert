package input

import (
	"github.com/zhjx922/alert/publisher"
	"sync"
)

type Keeper struct {
	inputs map[string]Input
	publisher *publisher.Publisher
}

func NewKeeper() *Keeper {
	return &Keeper{}
}

func (k *Keeper) SetPublisher(publisher *publisher.Publisher) {
	k.publisher = publisher
}

func (k *Keeper) Run(config *Config) error {

	var wg sync.WaitGroup

	// 初始化Input
	for _, v := range config.Inputs {
		in := NewInput(v)
		wg.Add(1)
		go func() {
			defer wg.Done()
			in.Run(k.publisher)
		}()
	}

	wg.Wait()

	return nil
}
