package demo_mocks

import "fmt"

// Q: how do you easily test a Method B that calls a Method A on the same struct?
// Answer below.

// Step 1: use function vars instead of methods directly.
// this seems like a strange abstraction layer at first,
// but its use will be apparent soon.
type Service struct {
	myFuncVar1 func(string) string
	myFuncVar2 func()
}

// -Impl suffix means the actual func implementation for normal program execution
func (s *Service) myMethod1Impl(input string) string {
	output := "This is my Method 1 being called! Input: " + input
	fmt.Println(output)
	return output
}

func (s *Service) myMethod2Impl() {
	s.myFuncVar1("Method 1 called from Method 2")
}

// Step 2: assign implementations to struct's function variables for normal operation.
// we use this constructor in non-test code.
func NewService() *Service {
	// setup
	newService := &Service{}

	// default, but these fields can be overridden as needed in tests.
	newService.myFuncVar1 = newService.myMethod1Impl
	newService.myFuncVar2 = newService.myMethod2Impl
	return newService
}

func NormalProgramFlow(myService *Service) {
	fmt.Println("Func calls on myService struct:")
	myService.myFuncVar1("first call") // myMethodImpl in normal code flow
	myService.myFuncVar2()             // myMethod2Impl in normal code flow
}
