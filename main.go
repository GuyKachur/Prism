package main

import (
	"fmt"
	"refract/server"
)

//modify power level color theme -
const (
	colorHolder  = "\033[38;5;%dm%s\033[39;49m "
	errorColor   = 203
	successColor = 49
)

// func PrintError(msg string) {
// 	PrintColor(errorColor, msg)
// }

func PrintColor(color int, msg string) {
	fmt.Printf(colorHolder, color, msg)
	fmt.Println()
}

func main() {
	PrintColor(successColor, "Starting server on port: 9090")
	server.NewServer()
	// for i := 1; i < 10; i++ {
	// 	fout, err := exec.Command("primitive", "-i", "/home/guy/projects/go/geo/prism/images/a.jpg", "-o", fmt.Sprintf("/home/guy/projects/go/geo/prism/images/output/a-shapes-%d.svg", i),
	// 		"-n", fmt.Sprint(i*50)).Output()
	// 	if err != nil {
	// 		l.Error(err)
	// 		return
	// 	}
	// 	l.Debug(fout)
	// }

	// con := refract.CreateDefault()
	// con.Input = "/home/guy/projects/go/geo/prism/images/a.jpg"
	// con.Output = "/home/guy/projects/go/geo/prism/images/a.svg"
	// con.Verbose = "on"
	// out, _, err := refract.Primitive(*con)
	// if err != nil {
	// 	l.Error(err)
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
