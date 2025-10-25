package greetings

import (
	"errors"
	"fmt"
	"math/rand"
)

func Hello(name string) (string, error) {
	if name == "" {
		return "", errors.New("empty name")
	}
	return fmt.Sprintf(randomFormat(), name), nil
}

// randomFormat return random greeting
func randomFormat() string {
	formats := []string{
		"Hi, %v. Welcome!",
		"Good day %v, how's it going",
		"Well met %v, my friend",
	}
	return formats[rand.Intn(len(formats))]
}
