package watchtower

import (
	"log"
	"time"
)

//New create new TowerWatch
func New(verbose ...bool) WatchTower {
	debug := func(v ...interface{}) {}
	debugf := func(format string, v ...interface{}) {}
	if len(verbose) > 0 && verbose[0] {
		debug = log.Println
		debugf = log.Printf
	}
	return &watcher{
		brokens:         map[string]Fixable{},
		fixingInProcess: new(AtomicBool),
		errmsgs:         map[string]string{},
		subs:            []func(){},
		fixInterval:     time.Second,
		debug:           debug,
		debugf:          debugf,
	}
}

//TowerWatch will watch if any infra is  unavailable, and will schedule to fix it
type WatchTower interface {
	AddWatchObject(fixables ...Fixable)
	AddInfraReadySubscribers(notifies ...func())
	IsBadInfrastructure() bool
	GetErrMessages() []string
	Run(chan bool)
}

//Fixable is object that will be watched by TowerWatch
type Fixable struct {
	Name    string
	Healthy func() bool
	Fix     func() error
	Err     string
}

type watcher struct {
	fixables        []Fixable
	brokens         map[string]Fixable
	fixingInProcess *AtomicBool
	errmsgs         map[string]string
	subs            []func()
	fixInterval     time.Duration
	debug           func(v ...interface{})                // log.Println
	debugf          func(format string, v ...interface{}) // log.Printf
}

func (w *watcher) AddWatchObject(fixables ...Fixable) {
	for _, f := range fixables {
		w.fixables = append(w.fixables, f)
	}
}

func (w *watcher) Run(done chan bool) {
	for{
		for _, f := range w.fixables {
			if !f.Healthy() {
				w.brokens[f.Name] = f
				w.errmsgs[f.Name] = f.Err
				w.fix()
			}
		}
		select {
		case <- done:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

func (w *watcher) AddInfraReadySubscribers(notifies ...func()) {
	for _, not := range notifies {
		if not != nil {
			w.subs = append(w.subs, not)
		}
	}
}

func (w *watcher) fix() {
	if len(w.brokens) == 0 || w.fixingInProcess.Get() {
		return
	}
	w.debug(">>> Start fixing")
	w.fixingInProcess.Set(true)
	go func() {
		defer func() {
			w.fixingInProcess.Set(false)
			w.debug(">>> Stop fixing")
		}()
		ticker := time.NewTicker(w.fixInterval)
		tick := ticker.C
		isStillFixing := new(AtomicBool)
		for {
			select {
			case <-tick:
				if isStillFixing.Get() {
					continue
				}
				isStillFixing.Set(true)
				for name, f := range w.brokens {
					w.debug(">>> Fixing ", name)
					err := f.Fix()
					if err != nil {
						w.debugf("Error fixing %s, got error : %+v\n", name, err)
					} else {
						delete(w.brokens, name)
						delete(w.errmsgs, name)
					}
				}
				if len(w.errmsgs) == 0 {
					ticker.Stop()
					w.notifSubscribers()
					return
				}
				isStillFixing.Set(false)
			}
		}
	}()
}

func (w *watcher) notifSubscribers() {
	if len(w.subs) == 0 {
		return
	}
	for _, notify := range w.subs {
		notify()
	}
}

func (w *watcher) IsBadInfrastructure() bool {
	return len(w.errmsgs) > 0 || w.fixingInProcess.Get()
}

func (w *watcher) GetErrMessages() []string {
	errs := []string{}
	for _, val := range w.errmsgs {
		errs = append(errs, val)
	}
	return errs
}
