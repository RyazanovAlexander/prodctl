//go:build mage

package main

import "fmt"

// Bundle builds bundle
func Bundle() error {
	fmt.Println("Done!")
	return nil
}

// Publish publishes artifacts
func Publish() error {
	fmt.Println("Done!")
	return nil
}

// Deploy deploys resources to the specified environment
// Params:
//   namespace: some description
func Deploy() error {
	fmt.Println("Done!")
	return nil
}

// Deletes resources
func Delete() error {
	fmt.Println("Done!")
	return nil
}
