
package main

import (
    "fmt"
    "crypto/rand"
)

func (t *SimpleChaincode) getUUID()([]byte, error){
     b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        fmt.Println("Error: ", err)
        return nil, err 
    }
    uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    byteArray := []byte(uuid)
    return byteArray, nil 
}