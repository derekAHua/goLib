package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

const defaultPipSize = 16

type (
	WalkFunc    func(item interface{}) (ret interface{}, err error) // WalkFunc 并发处理函数(具体业务执行逻辑)
	CollectFunc func(result interface{})                            // CollectFunc 收集结果

	ErrGroup interface {
		// Run
		// input:入参,必须是slice; do:并发处理函数,参数为input中元素; collect:收集结果; opts:可选项; pip:返回;
		Run(input interface{}, do WalkFunc, collect CollectFunc, opts ...ErrGroupOption)
		Error() error // wait and return error
		Go(f func() (err error))
	}

	defaultSafeErrGroup struct {
		ctx context.Context

		wg sync.WaitGroup

		errOnce sync.Once
		err     error
		cancel  func()
	}

	walkResult struct {
		result interface{}
		err    error
	}
)

func (g *defaultSafeErrGroup) Go(f func() error) {
	g.g(func() error {
		done := make(chan error)

		go func() {
			var err error
			defer func() {
				done <- err
				close(done)
			}()

			defer func() {
				if r := recover(); r != nil {
					err = errors.New(fmt.Sprintf("panic:[%v]", r))
				}
			}()
			err = f()
		}()

		select {
		case err := <-done:
			if err != nil {
				return err
			}

			return nil
		case <-g.ctx.Done():
			go func() { <-done }()
			return g.ctx.Err()
		}
	})
}

func (g *defaultSafeErrGroup) Run(input interface{}, do WalkFunc, collect CollectFunc, opts ...ErrGroupOption) {
	option := buildOptions(opts...)

	var pip <-chan interface{}
	if option.limit {
		pip = g.runLimited(input, do, option)
	} else {
		pip = g.runUnLimited(input, do, option)
	}

	for r := range pip {
		collect(r)
	}

	return
}

func (g *defaultSafeErrGroup) Error() error {
	g.wg.Wait()

	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

func (g *defaultSafeErrGroup) runUnLimited(input interface{}, do WalkFunc, option *rxOptions) (pip chan interface{}) {
	pip = make(chan interface{}, defaultPipSize)

	go func() {
		for v := range from(input) {
			g.do(do, v, pip, nil, false)
		}

		g.wg.Wait()
		close(pip)
	}()

	return
}

func (g *defaultSafeErrGroup) runLimited(input interface{}, do WalkFunc, option *rxOptions) (pip chan interface{}) {
	pip = make(chan interface{}, defaultPipSize)

	go func() {
		pool := make(chan struct{}, option.workers) // Limit by channel.
		for v := range from(input) {
			pool <- struct{}{}
			g.do(do, v, pip, pool, true)
		}

		g.wg.Wait()
		close(pip)
	}()

	return
}

func (g *defaultSafeErrGroup) do(w WalkFunc, item interface{}, pip chan interface{}, limitPool <-chan struct{}, limit bool) {
	g.g(func() error {
		if limit {
			defer func() { <-limitPool }()
		}

		done := make(chan walkResult)

		go func() {
			var ret walkResult

			defer func() {
				done <- ret
				close(done)
			}()

			defer func() {
				if r := recover(); r != nil {
					ret.err = errors.New(fmt.Sprintf("panic:[%v]", r))
				}
			}()
			ret.result, ret.err = w(item)
		}()

		select {
		case ret := <-done:
			if ret.err != nil {
				return ret.err
			}

			pip <- ret.result
			return nil
		case <-g.ctx.Done():
			go func() { <-done }()
			return g.ctx.Err()
		}
	})
}

func (g *defaultSafeErrGroup) g(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func NewErrGroup(ctx context.Context) ErrGroup {
	c, cancel := context.WithCancel(ctx)
	return &defaultSafeErrGroup{ctx: c, cancel: cancel}
}

// from: input must reflect.Slice.
func from(input interface{}) <-chan interface{} {
	switch reflect.TypeOf(input).Kind() {
	case reflect.Slice:
		obj := reflect.ValueOf(input)
		source := make(chan interface{})
		go func() {
			defer close(source)
			for i := 0; i < obj.Len(); i++ {
				source <- obj.Index(i).Interface()
			}
		}()
		return source
	default:
		panic("Using ErrGroup input must be slice.")
	}
}

type (
	ErrGroupOption func(opts *rxOptions)

	rxOptions struct {
		limit   bool // 是否控制协程数量
		workers int  // 协程数量
	}
)

func buildOptions(opts ...ErrGroupOption) *rxOptions {
	rx := new(rxOptions)
	for _, f := range opts {
		f(rx)
	}
	return rx
}

func WithWorkers(workerNum int) ErrGroupOption {
	if workerNum < 1 {
		panic("WorkerNum must more than zero.")
	}
	return func(opts *rxOptions) {
		opts.workers = workerNum
		opts.limit = true
	}
}
