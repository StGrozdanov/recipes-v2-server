package utils

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint used for printing structures such as structs, map in a human friendly way
func PrettyPrint(input interface{}) {
	bytes, _ := json.MarshalIndent(input, "", "")
	fmt.Println(string(bytes))
}
