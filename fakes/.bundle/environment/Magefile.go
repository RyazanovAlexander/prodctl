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

// Deploy deploys environment and release
// Params:
//   namespace: some description
func Deploy(namespace string) error {
	fmt.Println("Deploy done! Namespace: " + namespace)
	return nil
}

// Deletes resources
func Delete() error {
	fmt.Println("Done!")
	return nil
}
