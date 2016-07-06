package main

import (
	"fmt"
	"net/http"
	"bufio"
	"io/ioutil"
)

func main() {
	resp, err := http.Get("https://api.github.com/users/kinvolk/events")

	if err != nil {
		return
	}

	fmt.Printf("%s", resp)


}
