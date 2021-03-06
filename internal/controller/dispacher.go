package controller

import (
	"context"
	"fmt"
	"sync"

	"gofluentd/internal/monitor"
	"gofluentd/internal/tagfilters"
	"gofluentd/library"
	"gofluentd/library/log"

	utils "github.com/Laisky/go-utils"
	"github.com/Laisky/zap"
)

type DispatcherCfg struct {
	InChan             chan *library.FluentMsg
	TagPipeline        tagfilters.TagPipelineItf
	NFork, OutChanSize int
}

// Dispatcher dispatch messages by tag to different concator
type Dispatcher struct {
	*DispatcherCfg
	tag2Concator *sync.Map               // tag:msgchan
	tag2Counter  *sync.Map               // tag:counter
	tag2Cancel   *sync.Map               // tag:cancel
	outChan      chan *library.FluentMsg // skip concator, direct to producer
	counter      *utils.Counter
}

// NewDispatcher create new Dispatcher
func NewDispatcher(cfg *DispatcherCfg) *Dispatcher {
	d := &Dispatcher{
		DispatcherCfg: cfg,
		outChan:       make(chan *library.FluentMsg, cfg.OutChanSize),
		tag2Concator:  &sync.Map{},
		tag2Counter:   &sync.Map{},
		tag2Cancel:    &sync.Map{},
		counter:       utils.NewCounter(),
	}
	if err := d.valid(); err != nil {
		log.Logger.Panic("config invalid", zap.Error(err))
	}

	log.Logger.Info("create Dispatcher",
		zap.Int("n_fork", d.NFork),
		zap.Int("out_chan_size", d.OutChanSize),
	)
	return d
}

func (d *Dispatcher) valid() error {
	if d.NFork <= 0 {
		d.NFork = 4
		log.Logger.Info("reset n_fork", zap.Int("n_fork", d.NFork))
	}

	if d.OutChanSize <= 0 {
		d.OutChanSize = 1000
		log.Logger.Info("reset out_chan_size", zap.Int("out_chan_size", d.OutChanSize))
	}

	return nil
}

// Run dispacher to dispatch messages to different concators
func (d *Dispatcher) Run(ctx context.Context) {
	log.Logger.Info("run dispacher...")
	d.registerMonitor()
	lock := &sync.Mutex{}

	for i := 0; i < d.NFork; i++ {
		go func() {
			var (
				inChanForEachTagi interface{}
				inChanForEachTag  chan<- *library.FluentMsg
				ok                bool
				err               error
				counterI          interface{}
				msg               *library.FluentMsg
			)
			defer log.Logger.Info("dispatcher exist with msg", zap.String("msg", fmt.Sprint(msg)))

			// send each message to appropriate tagfilter by `tag`
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok = <-d.InChan:
					if !ok {
						log.Logger.Info("inchan closed")
						return
					}
				}

				d.counter.Count()
				if inChanForEachTagi, ok = d.tag2Concator.Load(msg.Tag); !ok {
					// create new inChanForEachTag
					lock.Lock()
					if inChanForEachTagi, ok = d.tag2Concator.Load(msg.Tag); !ok { // double check
						// new tag, create new tagfilter and its inchan
						log.Logger.Info("got new tag", zap.String("tag", msg.Tag))
						ctx2Tag, cancel := context.WithCancel(ctx)
						if inChanForEachTag, err = d.TagPipeline.Spawn(ctx2Tag, msg.Tag, d.outChan); err != nil {
							log.Logger.Error("try to spawn new tagpipeline got error",
								zap.Error(err),
								zap.String("tag", msg.Tag))
							cancel()
							continue
						} else {
							d.tag2Counter.Store(msg.Tag, utils.NewCounter())
							d.tag2Cancel.Store(msg.Tag, cancel)
							// tag2Concator should put after tag2Counter & tag2Cancel,
							// because the mutex only check whether tag2Concator has `msg.Tag`.
							d.tag2Concator.Store(msg.Tag, inChanForEachTag)
							go func(tag string) {
								<-ctx2Tag.Done()
								log.Logger.Info("remove tag in dispatcher", zap.String("tag", tag))
								lock.Lock()
								d.tag2Concator.Delete(tag)
								d.tag2Counter.Delete(tag)
								d.tag2Cancel.Delete(tag)
								lock.Unlock()
							}(msg.Tag)
						}
					} else {
						inChanForEachTag = inChanForEachTagi.(chan<- *library.FluentMsg)
					}

					lock.Unlock()
				} else {
					inChanForEachTag = inChanForEachTagi.(chan<- *library.FluentMsg)
				}

				// count
				if counterI, ok = d.tag2Counter.Load(msg.Tag); !ok {
					log.Logger.Panic("counter must exists", zap.String("tag", msg.Tag))
				}
				counterI.(*utils.Counter).Count()

				// put msg into tagfilter's inchan
				select {
				case inChanForEachTag <- msg:
				default:
					log.Logger.Warn("discard msg since tagfilter's inchan is blocked", zap.String("tag", msg.Tag))
				}
			}
		}()
	}
}

func (d *Dispatcher) registerMonitor() {
	monitor.AddMetric("dispatcher", func() map[string]interface{} {
		metrics := map[string]interface{}{
			"msgPerSec": d.counter.GetSpeed(),
			"msgTotal":  d.counter.Get(),
			"config": map[string]interface{}{
				"n_fork":        d.NFork,
				"out_chan_size": d.OutChanSize,
			},
		}
		d.tag2Counter.Range(func(tagi interface{}, ci interface{}) bool {
			metrics[tagi.(string)+".MsgPerSec"] = ci.(*utils.Counter).GetSpeed()
			metrics[tagi.(string)+".msgTotal"] = ci.(*utils.Counter).Get()
			return true
		})

		d.tag2Concator.Range(func(tagi interface{}, ci interface{}) bool {
			metrics[tagi.(string)+".ChanLen"] = len(ci.(chan<- *library.FluentMsg))
			metrics[tagi.(string)+".ChanCap"] = cap(ci.(chan<- *library.FluentMsg))
			return true
		})
		return metrics
	})
}

func (d *Dispatcher) GetOutChan() chan *library.FluentMsg {
	return d.outChan
}
