/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"math/rand"
	"time"

	"github.com/luizbafilho/mokuro-vox/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
