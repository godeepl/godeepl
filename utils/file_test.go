package utils

import (
    "fmt"
    "os"
    "strings"
    "testing"
)

func TestIsFileExist(t *testing.T) {
    fmt.Println(IsFileExist("/path"))
    fmt.Println(IsFileExist("/"))
    fmt.Println(IsFileExist("/usr/local/opt/go/libexec/bin/go"))
}

func TestReadMarkdown(t *testing.T) {
    pwd, _ := os.Getwd()
    filename := fmt.Sprintf("%s/samples/test-1.md", pwd)
    result, err := ReadMarkdown(filename)
    if err != nil {
        fmt.Println("error:", err.Error())
    }

    fmt.Println(fmt.Sprintf("%v", result))

    // write to markdown
    newfilename := strings.Replace(filename, ".md", ".deepl.md", 1)
    err = WriteMarkdown(newfilename, result)
    if err != nil {
        return
    }
}
