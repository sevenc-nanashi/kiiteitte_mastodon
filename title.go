package main

import (
	"fmt"

	"github.com/fatih/color"
)

func Title() {
	colorKiite := color.New(color.Bold).SprintFunc()
	colorI := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	colorTte := color.New(color.FgHiMagenta, color.Bold).SprintFunc()
	fmt.Printf("%s%s%s  \n", colorKiite("Kiite"), colorI("i"), colorTte("tte"))
	fmt.Println("  https://github.com/sevenc-nanashi/kiiteitte_vocalodon\n")

}
