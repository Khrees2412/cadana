package util

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(p interface{}) {
	pretty, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(pretty))
}
