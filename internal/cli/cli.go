package cli

import "fmt"

var (
	success string = "\x1b[32m\u2713\x1b[0m"
	failure string = "\x1b[31m\u2717\x1b[0m"
)

func Success() { fmt.Printf(" %s\n", success) }
func Failure() { fmt.Printf(" %s\n", failure) }

func PrintStatus(failed bool) {
	if failed {
		Failure()
	} else {
		Success()
	}
}
