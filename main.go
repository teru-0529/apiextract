/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package main

import (
	"fmt"

	"github.com/teru-0529/apiextract/cmd"
)

const version = "0.0.0"

func main() {
	fmt.Printf("version: %s\n", version)
	cmd.Execute()
}
