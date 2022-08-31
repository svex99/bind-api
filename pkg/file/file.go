package file

import (
	"io/ioutil"
	"os"
	"strings"
)

// Replaces content of the file that matches `old` string with the `new` one.
// If backup is true the old file is preserved with .bak appended to its name.
func ReplaceContent(filePath, old, new string, backup bool) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	updatedContent := strings.Replace(string(content), old, new, 1)

	if backup {
		if err := os.Rename(filePath, filePath+".bak"); err != nil {
			return err
		}
	}

	if err := os.WriteFile(filePath, []byte(updatedContent), 0666); err != nil {
		return err
	}

	return nil
}

func AddContent(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}