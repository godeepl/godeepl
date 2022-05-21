package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func IsFileExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}

		if os.IsNotExist(err) {
			return false
		}

		return false
	} else {
		return !fileInfo.IsDir()
	}
}

func Write(filename string, data []byte, fileMode os.FileMode) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	defer file.Close()
	n, err := file.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	return err
}

type FileType string

const (
	Text FileType = "Text"
	Code FileType = "Code"
)

type ContentItem struct {
	Content []string `json:"content"`
	Type    FileType `json:"type"`
}

// ReadMarkdown read text from markdown file, and return an array
func ReadMarkdown(filename string) ([]ContentItem, error) {
	var result []ContentItem

	file, err := os.OpenFile(filename, os.O_RDONLY, 0664)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	startMD := false
	codeContentItem := ContentItem{
		Content: []string{},
		Type:    Code,
	}
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		str := string(line)
		if strings.HasPrefix(str, "```") && startMD == false {
			startMD = true
			codeContentItem.Content = append(codeContentItem.Content, str)
		} else if strings.HasPrefix(str, "```") && startMD == true {
			startMD = false
			codeContentItem.Content = append(codeContentItem.Content, str)
			result = append(result, codeContentItem)
            codeContentItem.Content = []string{}
		} else if startMD == true {
			codeContentItem.Content = append(codeContentItem.Content, str)
		} else {
			result = append(result, ContentItem{
				Content: []string{string(line)},
				Type:    Text,
			})
		}

	}

	return result, nil
}

// WriteMarkdown write markdown content to file
func WriteMarkdown(filename string, items []ContentItem) error {
	fmt.Println("begin to write result to", filename)
	var file *os.File
	var err error
	if IsFileExist(filename) {
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	var content string
	for _, item := range items {
		if content == "" {
			content = strings.Join(item.Content, "\n")
		} else {
			content = content + "\n" + strings.Join(item.Content, "\n")
		}
	}

	w := bufio.NewWriter(file)
	n, err := w.WriteString(content)

	fmt.Printf("write %d bytes.", n)
	err = w.Flush()
	if err != nil {
		fmt.Println("write to file", filename, "error", err.Error())
		os.Exit(1)
	}

	return nil
}
