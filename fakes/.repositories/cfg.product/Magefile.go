//go:build mage

package main

import "fmt"

// Deploy deploys environment and release
// Params:
//   namespace: some description
func Deploy(namespace string) error {
	fmt.Println("Deploy done! Namespace: " + namespace)
	return nil
}

// Delete deletes environment and release
func Delete() error {
	fmt.Println("Delete done!")
	return nil
}

// Ready checks the product for readiness to perform tasks
func Ready() error {
	fmt.Println("Ready done!")
	return nil
}

// Binaries uploads images and charts to the specified container registry
// Params:
//   containerRegistry: some description
func Binaries(containerRegistry string) error {
	fmt.Println("Images done!")
	return nil
}
