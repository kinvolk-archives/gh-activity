package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// get the data
	resp, err := http.Get("https://api.github.com/users/kinvolk/events")
	if err != nil {
		fmt.Printf("error getting the data from github", err)
		return
	}

	// output data to the file
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err = ioutil.WriteFile("myfile.json", buf.Bytes(), 0666)
	if err != nil {
		return
	}
}
