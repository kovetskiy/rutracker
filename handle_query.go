package main

import "fmt"

func handleQuery(args map[string]interface{}, config *Config) error {
	var (
		query = args["--query"].(string)
	)

	fmt.Printf("XXXXXX handle_query.go:7 query: %#v\n", query)

	return nil
}
