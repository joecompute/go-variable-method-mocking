# Method-Level Mocking in Go Using Function Variables

This repository demonstrates a technique for mocking methods on a single Go struct without using interfaces. This is particularly useful when you need to test a method that calls another method on the same struct.

## The Problem

How do you easily test a method `B` that calls a method `A` on the same struct?

```go
type Service struct{}

func (s *Service) MethodA() {
    // ... some logic
}

func (s *Service) MethodB() {
    s.MethodA() // How do we mock this call in a test for MethodB?
}
```

Directly testing `MethodB` in isolation is difficult because it has a hard dependency on `MethodA`.

## The Technique: Function Variables for Dependency Injection

This approach decouples the method's implementation from its invocation by using function variables within the struct.

### Step 1: Use Function Variables in the Struct

Instead of methods, define fields in your struct that are function types.

```go
// From: demo_mocks/demo_mocks.go

type Service struct {
	myFuncVar1 func(string) string
	myFuncVar2 func()
}
```

The actual logic is defined in separate implementation functions (e.g., `myMethod1Impl`).

```go
// From: demo_mocks/demo_mocks.go

func (s *Service) myMethod1Impl(input string) string {
	output := "This is my Method 1 being called! Input: " + input
	fmt.Println(output)
	return output
}

func (s *Service) myMethod2Impl() {
	s.myFuncVar1("Method 1 called from Method 2")
}
```

### Step 2: Assign Implementations in a Constructor

For normal program operation, a constructor assigns the real implementation functions to the corresponding function variables.

```go
// From: demo_mocks/demo_mocks.go

func NewService() *Service {
	newService := &Service{}

	// Assign the "real" methods to the function variables for normal use.
	newService.myFuncVar1 = newService.myMethod1Impl
	newService.myFuncVar2 = newService.myMethod2Impl
	return newService
}
```

### Step 3: Override Function Variables in Tests

In your tests, you can easily replace the real implementations with mocks. This allows you to isolate the function under test.

#### Simple Mock (No Framework)

For simple cases, you can use an anonymous function to track calls or control return values.

```go
// From: demo_mocks/demo_mocks_test.go

func TestBasic(t *testing.T) {
	myService := NewService()
    callCount := 0

	// Override the function variable with a simple mock
	myService.myFuncVar1 = func(input string) string {
		callCount++
		return "mocked return"
	}

	myService.myFuncVar2() // This will now call our mock instead of the real implementation

	assert.Equal(t, 1, callCount)
}
```

#### Integration with Mocking Libraries (testify/mock)

For more advanced testing, this technique integrates seamlessly with mocking libraries like `testify/mock`.

```go
// From: demo_mocks/demo_mocks_test.go

func (m *ServiceMock) myMethod1Instrumented(input string) string {
	args := m.Called(input)
	return args.String(0)
}

func TestWTestifyMockBasic(t *testing.T) {
	wrappedWithMock := ServiceMock{Service: NewService()}
    // Assign the instrumented mock method
	wrappedWithMock.myFuncVar1 = wrappedWithMock.myMethod1Instrumented

    // Set up expectations
	wrappedWithMock.On("myMethod1Instrumented", "Method 1 called from Method 2").Return("my custom return").Once()

	// Run the code that calls the function
    wrappedWithMock.myFuncVar2()

	// Assert that the expectations were met
	wrappedWithMock.AssertExpectations(t)
}
```

## Benefits

*   **No Interfaces Needed**: Avoids the overhead of creating interfaces just for the sake of mocking a single struct's methods.
*   **Improved Readability**: The intent of the code and its tests becomes clearer.
*   **Flexible**: Easily swap implementations in and out. You can use simple anonymous functions for basic mocks or integrate with powerful mocking frameworks for more complex scenarios.
*   **Isolated Testing**: Enables true unit testing of methods by removing internal dependencies within the same struct.
