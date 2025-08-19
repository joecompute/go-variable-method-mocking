package main

import "joetesting/demo_mocks"

func main() {
	myService := demo_mocks.NewService()
	demo_mocks.NormalProgramFlow(myService)
}
