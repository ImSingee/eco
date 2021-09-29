package eco

import (
	"fmt"
	"strings"
)

type Collector struct {
	Errors []error
}

func New() Collector {
	return Collector{}
}

func (e Collector) Err() error {
	switch len(e.Errors) {
	case 0:
		return nil
	case 1:
		return e.Errors[0]
	default:
		return e
	}
}

func (e Collector) Error() string {
	switch len(e.Errors) {
	case 0:
		return ""
	case 1:
		return e.Errors[0].Error()
	default:
		ss := make([]string, len(e.Errors))

		for i, err := range e.Errors {
			ss[i] = err.Error()
		}

		return fmt.Sprintf("%d errors: %s", len(e.Errors), strings.Join(ss, ", "))
	}
}

func (e Collector) Unwrap() error {
	switch len(e.Errors) {
	case 0:
		return nil
	default:
		return e.Errors[0]
	}
}

func (e Collector) IsError() bool {
	return len(e.Errors) != 0
}

func (e *Collector) Collect(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// C is same as Collect, but also return received error
func (e *Collector) C(err error) error {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}

	return err
}

func (e *Collector) Do(f func() error) {
	e.Collect(e.do(f))
}

// D is same as Do, but also return received error
func (e *Collector) D(f func() error) (err error) {
	return e.C(e.do(f))
}

func (e *Collector) do(f func() error) (err error) {
	defer func() {
		v := recover()
		if v != nil {
			err = fmt.Errorf("panic: %v", v)
		}
	}()

	err = f()
	return
}

// Process do some functions but only if previous is not error
func (e *Collector) Process(processes ...func() error) error {
	for _, p := range processes {
		if e.IsError() {
			return e.Err()
		}

		e.Do(p)
	}

	return e.Err()
}

// Do some things until any error
func Do(processes ...func() error) error {
	e := New()
	return e.Process(processes...)
}
