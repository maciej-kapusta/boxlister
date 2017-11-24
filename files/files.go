package files

import (
	"io/ioutil"
	"os"
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func FillGenerated(fileName *string, buffer bytes.Buffer) {
	content, e := ioutil.ReadFile(*fileName)
	checkError(e)
	file, e := os.Create(*fileName)
	checkError(e)
	defer file.Close()
	scanner := bufio.NewScanner(bytes.NewReader(content))
	generating := false
	for scanner.Scan() {
		line := scanner.Text()
		if ! generating {
			_, e := fmt.Fprintln(file, line)
			checkError(e)
		}

		if strings.Index(line, "#generated") == 0 {
			generating = true
			_, e := fmt.Fprint(file, buffer.String())

			checkError(e)
			break
		}
	}
	if !generating {
		fmt.Println("Not generated.")
	}
}
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
