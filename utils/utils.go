package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrintStruct(obj interface{}) {
	objJson, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(objJson))
}
