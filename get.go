package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		var results []string

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		results = append(results, string(body))
		fmt.Println(results)
	})
	http.ListenAndServe(":8182", nil)
}
