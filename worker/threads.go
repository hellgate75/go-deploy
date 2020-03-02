package worker

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Runnable interface {
	Run() error
	Stop() error
	Kill() error
	Pause() error
	Resume() error
	IsRunning() bool
	IsComplete() bool
	UUID() string
}

type ThreadPool interface {
	Add(r Runnable) error
	Start() error
	Stop() error
	Reset() error
	IsStarted() bool
	IsComplete() bool
	WaitFor() error
	SetErrorHandler(h ThreadErrorHandler)
}

type ThreadErrorHandler interface {
	HandleError(e error)
}

//runtime.GOMAXPROCS(runtime.NumCPU())
type threadPool struct {
	sync.RWMutex
	maxThreads int64
	parallel   bool
	threads    []Runnable
	running    bool
	errHandler ThreadErrorHandler
}

func (tp *threadPool) SetErrorHandler(h ThreadErrorHandler) {
	tp.errHandler = h
}

func (tp *threadPool) WaitFor() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	for !tp.IsComplete() {
		time.Sleep(10 * time.Second)
	}
	return err
}

func (tp *threadPool) IsComplete() bool {
	var result bool = true
	defer func() {
		if r := recover(); r != nil {
			result = false
		}
		tp.RUnlock()
	}()
	tp.RLock()
	for _, r := range tp.threads {
		if r.IsRunning() || !r.IsComplete() {
			result = false
			break
		}
	}
	return result
}

func (tp *threadPool) IsStarted() bool {
	return tp.running
}

func (tp *threadPool) Reset() error {
	if tp.running {
		return errors.New("Unable to stop running ThreadPool")
	}
	return nil
}

func (tp *threadPool) Stop() error {
	tp.running = false
	return nil
}

func (tp *threadPool) Start() error {
	tp.running = true
	return nil
}

func (tp *threadPool) clean() {
	var isLocked bool = false
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("worker.ThreadPool Failure cleaning pool -> Details: %v", r)
		}
		if isLocked {
			tp.Unlock()
		}
	}()
	for idx, t := range tp.threads {
		if !t.IsRunning() || t.IsComplete() {
			t.Kill()
			isLocked = true
			tp.Lock()
			if idx == 0 {
				//Removing first item of the list
				if len(tp.threads) > 1 {
					tp.threads = tp.threads[1:]
				} else {
					tp.threads = make([]Runnable, 0)
				}
			} else if idx == len(tp.threads)-1 {
				//Removing last item of the list
				if len(tp.threads) > 1 {
					tp.threads = tp.threads[:len(tp.threads)-1]
				} else {
					tp.threads = make([]Runnable, 0)
				}
			} else {
				//Removing an item into the list except boundaries
				lower := tp.threads[:idx]
				upper := tp.threads[idx+1:]
				tp.threads = lower
				tp.threads = append(tp.threads, upper...)
			}
			tp.Unlock()
			isLocked = false
		}
		runtime.GC()
	}
}

func (tp *threadPool) run() {
	var isLocked bool = false
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("worker.ThreadPool Failure running pool -> Details: %v", r)
		}
		if isLocked {
			tp.Unlock()
		}
	}()
	if tp.parallel {
		if tp.maxThreads == 0 {
			//threads all together
			for !tp.IsComplete() {
				var makedForClean bool = false
				for _, t := range tp.threads {
					if !t.IsRunning() {
						if !t.IsComplete() {
							go func() {
								err := t.Run()
								if err != nil && tp.errHandler != nil {
									tp.errHandler.HandleError(err)
								}
								t.Kill()
							}()
						} else {
							makedForClean = true
						}
					}
				}
				if makedForClean {
					go tp.clean()
				}
				time.Sleep(5 * time.Second)
			}
		} else {
			//run max thread together and remove completed
			for !tp.IsComplete() {
				var makedForClean bool = false
				var numActive int64 = 0
				for _, t := range tp.threads {
					if numActive < tp.maxThreads {
						if !t.IsRunning() {
							if !t.IsComplete() {
								go func() {
									err := t.Run()
									if err != nil && tp.errHandler != nil {
										tp.errHandler.HandleError(err)
									}
									t.Kill()
								}()
							} else {
								makedForClean = true
							}
						} else {
							numActive += 1
						}
					}
				}
				if makedForClean {
					go tp.clean()
				}
				time.Sleep(5 * time.Second)
			}

		}
	} else {
		for !tp.IsComplete() {
			if len(tp.threads) > 0 {
				if tp.threads[0] == nil {
					if len(tp.threads) > 1 {
						isLocked = true
						tp.Lock()
						tp.threads = tp.threads[1:]
					} else {
						tp.threads = make([]Runnable, 0)
					}
					continue
				}
				if !tp.threads[0].IsRunning() {
					if !tp.threads[0].IsComplete() {
						go func() {
							err := tp.threads[0].Run()
							if err != nil && tp.errHandler != nil {
								tp.errHandler.HandleError(err)
							}
							tp.threads[0].Kill()
						}()
					} else {
						if len(tp.threads) > 1 {
							isLocked = true
							tp.Lock()
							tp.threads = tp.threads[1:]
						} else {
							tp.threads = make([]Runnable, 0)
						}
					}
				}
			} else {
				tp.running = false
			}
			runtime.GC()
			time.Sleep(5 * time.Second)
		}
	}
}

func (tp *threadPool) Add(r Runnable) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		tp.Unlock()
	}()
	tp.Lock()
	tp.threads = append(tp.threads, r)
	return err
}

func NewThreadPool(maxThreads int64, parallel bool) ThreadPool {
	return &threadPool{
		maxThreads: maxThreads,
		parallel:   parallel,
		threads:    make([]Runnable, 0),
	}
}
