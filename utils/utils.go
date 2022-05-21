package utils

import (
    "fmt"
    "math/rand"
    "time"
)

func Sleep() {
    rand.Seed(time.Now().Unix())
    num := rand.Int31n(20)
    if num < 5 {
        num += 5
    }
    fmt.Println("randam sleep ", num, "s")
    time.Sleep(time.Duration(num) * time.Second)
}
