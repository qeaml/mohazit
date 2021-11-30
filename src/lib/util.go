package lib

import "fmt"

func moreArgs(txt string) error {
	return fmt.Errorf("not enough arguments: %s", txt)
}

func badType(txt string) error {
	return fmt.Errorf("wrong type: %s", txt)
}

func badState(txt string) error {
	return fmt.Errorf("invalid state: %s", txt)
}
