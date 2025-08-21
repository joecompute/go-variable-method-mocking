package demo_mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Benefit: testing without need for more heavyweight interfaces
// code is more readable and it is much easier to swap methods in and out.

type ServiceMock struct {
	*Service
	mock.Mock
	count           int
	capturedReturns []string
}

// basic use of the wrapped funcs via func variables/mock
// no need to use testify/mock in many cases
func TestBasic(t *testing.T) {
	wrappedWithMock := ServiceMock{Service: NewService()}

	// simple mock: basic count w arg passthru
	wrappedWithMock.myFuncVar1 = func(input string) string {
		wrappedWithMock.count++ // now we can count func calls!
		return input            // custom return, e.g., in the case, func becomes as passthrough
	}

	NormalProgramFlow(wrappedWithMock.Service)

	// test result: the first function was called twice per normalProgramFlow
	assert.Equal(t, 2, wrappedWithMock.count)

	// SCENARIO 2: basic count injection, plus preserving old method call
	// (not recommended to do often, but possible if really needed)
	wrappedWithMock.myFuncVar1 = func(input string) string {
		wrappedWithMock.count++                     // still counting func calls
		return wrappedWithMock.myMethod1Impl(input) // call the old method and return the result
	}

	NormalProgramFlow(wrappedWithMock.Service)

	assert.Equal(t, 4, wrappedWithMock.count)
}

func (m *ServiceMock) myMethod1Instrumented(input string) string {
	args := m.Called(input)
	return args.String(0)
}

// leverage powerful testify/mock instrumentation via instrumented method above
func TestWTestifyMockBasic(t *testing.T) {
	wrappedWithMock := ServiceMock{Service: NewService()}
	wrappedWithMock.myFuncVar1 = wrappedWithMock.myMethod1Instrumented

	wrappedWithMock.On("myMethod1Instrumented", "first call").Return("my custom return").Once()
	wrappedWithMock.On("myMethod1Instrumented",
		"Method 1 called from Method 2").Return("my custom return 2").Once()

	NormalProgramFlow(wrappedWithMock.Service)

	// use testify/mock func assertions
	wrappedWithMock.AssertExpectations(t)

}

// not recommended, but possible when necessary
// (e.g. when you have to call an API for an integration test)
func (m *ServiceMock) myMethod1InstrumentedWithOrigCallTracking(input string) string {
	args := m.Called(input)
	originalReturn := m.Service.myMethod1Impl(input)
	m.capturedReturns = append(m.capturedReturns, originalReturn)
	return args.String(0)
}

func TestWTestifyMockPassthrough(t *testing.T) {
	wrappedWithMock := ServiceMock{Service: NewService()}
	// take it a step further and let the original call happen, too
	wrappedWithMock.myFuncVar1 = wrappedWithMock.myMethod1InstrumentedWithOrigCallTracking
	wrappedWithMock.On("myMethod1InstrumentedWithOrigCallTracking",
		"first call").Return("my custom return").Once()
	wrappedWithMock.On("myMethod1InstrumentedWithOrigCallTracking",
		"Method 1 called from Method 2").Return("my custom return 2").Once()
	NormalProgramFlow(wrappedWithMock.Service)
	assert.Len(t, wrappedWithMock.capturedReturns, 2)
	wrappedWithMock.AssertExpectations(t)
}
