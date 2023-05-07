package main

import "log"

func main() {
	_, err := readData()

	if err != nil {
		log.Fatal(err)
	}
}
