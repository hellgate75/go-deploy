package pool

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Defines a single running unit of operational code
type Runnable interface {
	// Executes Runnable code
	Run() error
	// Gracefully stops Runnable code
	Stop() error
	// Forces to interrupt the thread
	Kill() error
	// Pause the execution if suits code purposes
	Pause() error
	// Resume the execution if suits code purposes
	Resume() error
	// Verify if component is running if suits code purposes
	IsRunning() bool
	// Verify if component is pause
	IsPaused() bool
	// Verify if code execution is consumed successfulyy or not
	IsComplete() bool
	// Retrieve unique Id of the component instance
	UUID() string
	// Returns the code execution duration until the request
	UpTime() time.Duration
}

// Defined Operational interface of a Thread Pool Manager
type ThreadPool interface {
	// Add new Runnable component in the ThreadPool
	Schedule(r Runnable) error
	// Start execution of ThreadPool
	Start() error
	// Stop gracefully execution of ThreadPool
	Stop() error
	// Pause temporarly execution of ThreadPool
	Pause() error
	// Resume paused execution of ThreadPool
	Resume() error
	//Reset state if ThreadPool is stopped and complete
	Reset() error
	// Verify is ThreadPool is running
	IsStarted() bool
	// Verify if all allocated threads are complete
	IsComplete() bool
	// Verify if ThreadPool is paused
	IsPaused() bool
	// Stop main thread waiting for ThreadPool completion
	WaitFor() error
	// Set Runnable errors listener, used to report errors occured
	// during ThreadPool operational time, executing scheduled
	// Runnable code
	SetErrorHandler(h ThreadErrorHandler)
	// Prints state of running processes and number of elements in the Queue
	State() string
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
	_size      int64
	_paused    bool
}

func format(format string, value interface{}, length int) string {
	var out = fmt.Sprintf(format, value)
	for len(out) < length {
		out += " "
	}
	if len(out) > length {
		out = out[:length-3] + "..."
	}
	return out
}

var (
	typeLen int = 15
	uuidLen int = 15
)

func (tp *threadPool) State() string {
	var out string = "Thread Pool Manager state:\n"
	out += "----------------------------------------------------------------------------\n"
	var complete, running, paused, waiting int64

	if len(tp.threads) > 0 {
		out += fmt.Sprintf("STATE     %s   %s   %s\n", format("%s", "TYPE", typeLen), format("%s", "UUID", uuidLen), "TIME")
	} else {
		out += "No threads scheduled\n"

	}
	for _, v := range tp.threads {
		if v.IsRunning() && !v.IsComplete() && !v.IsPaused() {
			running += 1
			out += fmt.Sprintf("Running   %s   %s   Up since %s\n", format("%T", v, typeLen), format("%s", v.UUID(), uuidLen), v.UpTime().String())
		} else if v.IsComplete() {
			complete += 1
			out += fmt.Sprintf("Complete  %s   %s   Done in %s\n", format("%T", v, typeLen), format("%s", v.UUID(), uuidLen), v.UpTime().String())
		} else if !v.IsPaused() {
			waiting += 1
			out += fmt.Sprintf("Waiting   %s   %s\n", format("%T", v, typeLen), format("%s", v.UUID(), uuidLen))
		} else {
			paused += 1
			out += fmt.Sprintf("Paused    %s   %s   Running for %s\n", format("%T", v, typeLen), format("%s", v.UUID(), uuidLen), v.UpTime().String())
		}
	}
	out += "----------------------------------------------------------------------------\n"
	out += fmt.Sprintf(" total elements: %v, complete: %v, running: %v, paused: %v, waiting: %v\n", tp._size, complete, running, paused, waiting)
	out += fmt.Sprintf(" ready: %v, complete: %v, parallel: %v, max threads: %v\n", tp.running, tp.IsComplete(), tp.parallel, tp.maxThreads)
	return out
}

func (tp *threadPool) IsPaused() bool {
	return tp._paused
}

func (tp *threadPool) Pause() error {
	for _, v := range tp.threads {
		if v.IsRunning() && !v.IsComplete() {
			err := v.Pause()
			if err != nil {
				return err
			}
		}
	}
	tp._paused = true
	return nil
}

func (tp *threadPool) Resume() error {
	for _, v := range tp.threads {
		if v.IsPaused() {
			err := v.Resume()
			if err != nil {
				return err
			}
		}
	}
	tp._paused = false
	return nil
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
		return errors.New("Unable to reset running ThreadPool")
	}
	if !tp.IsComplete() {
		return errors.New("Unable to reset uncomplete ThreadPool, please wait threads finish the work")
	}
	tp.threads = make([]Runnable, 0)
	tp._size = 0
	return nil
}

func (tp *threadPool) Stop() error {
	tp.running = false
	return nil
}

func (tp *threadPool) Start() error {
	tp.running = true
	tp._paused = false
	go tp.run()
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
								defer func() {
									if r := recover(); r != nil {
										err := errors.New(fmt.Sprintf("%v", r))
										if err != nil && tp.errHandler != nil {
											tp.errHandler.HandleError(err)
											t.Kill()
										}
									}
								}()
								err := t.Run()
								if err != nil && tp.errHandler != nil {
									tp.errHandler.HandleError(err)
									t.Kill()
								}
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
										t.Kill()
									}
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
						t := tp.threads[0]
						go func() {
							defer func() {
								if r := recover(); r != nil {
									err := errors.New(fmt.Sprintf("%v", r))
									if err != nil && tp.errHandler != nil {
										tp.errHandler.HandleError(err)
										t.Kill()
									}
								}
								t.Stop()
							}()
							err := t.Run()
							if err != nil && tp.errHandler != nil {
								tp.errHandler.HandleError(err)
								t.Kill()
							}
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

func (tp *threadPool) Schedule(r Runnable) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		tp.Unlock()
	}()
	tp.Lock()
	tp.threads = append(tp.threads, r)
	tp._size += 1
	return err
}

func NewThreadPool(maxThreads int64, parallel bool) ThreadPool {
	return &threadPool{
		maxThreads: maxThreads,
		parallel:   parallel,
		threads:    make([]Runnable, 0),
	}
}
