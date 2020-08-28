package main

import (
	"fmt"
	"os/exec"
)

//modify power level color theme -
const (
	colorHolder  = "\033[38;5;%dm%s\033[39;49m "
	errorColor   = 203
	successColor = 49
)

func PrintError(msg string) {
	PrintColor(errorColor, msg)
}

func PrintColor(color int, msg string) {
	fmt.Printf(colorHolder, color, msg)
	fmt.Println()
}

func main() {
	for i := 0; i < 5; i++ {
		fout, err := exec.Command("primitive", "-i", "/home/guy/projects/go/geo/middleware/images/a.jpg", "-o", fmt.Sprintf("/home/guy/projects/go/geo/middleware/images/a-shapes-%d.svg", i),
			"-n", fmt.Sprint(i*50)).Output()
		if err != nil {
			PrintError(err.Error())
			return
		}
		fmt.Println(fout)
	}

	// con := api.CreateDefault()
	// con.Input = "/home/guy/projects/go/geo/middleware/images/a.jpg"
	// con.Output = "/home/guy/projects/go/geo/middleware/images/a.svg"
	// con.Verbose = "on"
	// out, err := api.Primitive(*con)
	// if err != nil {
	// 	PrintError(err.Error())
	// 	return
	// }
	// PrintColor(successColor,
	// 	fmt.Sprintf("Sucess! config: %s", con.CommandForm()))
	// fmt.Println(string(out))
}

// func main() {
// 	for j := 0; j < 256; j++ {
// 		fmt.Printf(colorHolder, j, fmt.Sprintf("%d", j))
// 		PrintColor(j, "Hello!")
// 	}
// }
