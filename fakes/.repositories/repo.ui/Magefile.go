//go:build mage

package main

import "fmt"

// Build builds source code
func Build() error {
	fmt.Println("Done!")
	return nil
}

// Bundle builds bundle
func Bundle() error {
	fmt.Println("Done!")
	return nil
}

// Test run tests
func Test() error {
	fmt.Println("Done!")
	return nil
}

// Scan run scan
func Scan() error {
	fmt.Println("Done!")
	return nil
}

// Publish publishes artifacts
func Publish() error {
	fmt.Println("Done!")
	return nil
}

// Deploy deploys resources to the specified environment
func Deploy() error {
	fmt.Println("Done!")
	return nil
}

// Deletes resources
func Delete() error {
	fmt.Println("Done!")
	return nil
}

// To run the pipeline locally, uncomment the following block below and replace the first line with "///go:build mage".
// func main() {
// 	BuildDeployTestRemovePipeline()
// }

// func BuildDeployTestRemovePipeline() {
// 	Build()
// 	Bundle()
// 	Scan()
// 	Publish()
// 	Deploy()
// 	Test()
// 	Delete()
// }
