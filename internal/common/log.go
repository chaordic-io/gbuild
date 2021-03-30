package common

import "fmt"

// We provide a Log interface, so we can achieve two goals:
// * Make sure our tests are not noisy
// * Potentiall introduce file logging at a later time as a configurable option
type Log interface {
	Println(a ...interface{}) (n int, err error)
	Printf(format string, a ...interface{}) (n int, err error)
}

type OSLog struct{}

type NoLog struct{}

func (log OSLog) Println(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func (log OSLog) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}

func (log NoLog) Println(a ...interface{}) (n int, err error) {
	return 0, nil
}

func (log NoLog) Printf(format string, a ...interface{}) (n int, err error) {
	return 0, nil
}
