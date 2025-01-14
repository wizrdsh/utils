package utils 

import (
	"bufio"
	"fmt"
	"strings"
)
// makes it to enable / disable for debugging, link this to .env or an argument and you would be golden.
var RedactedConfig = struct {
	Enabled bool
	Phrases []string
	Message string
}{
	Enabled: true,
	Phrases: []string{"hide_this", "password123"},
	Message: "REDACTED",
}


// Redact(string)
// Redacts any string found within a string (or files if converted to string)
// This is useful for password redactions in plain text outputs or stdout (print, format, etc) statements.
// Accepts a list of values to redact inside
// Usage: fmt.Println(Redact(supeSecretOutputFilledWithLotsOfSecrets))

func Redact(output string) string {
	var results strings.Builder

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		scanline := fmt.Sprint(scanner.Text(), "\n")
		if RedactedConfig.Enabled {
			for _, value := range RedactedConfig.Phrases {
				scanline = strings.ReplaceAll(scanline, value, RedactedConfig.Message)
			}
		}
		results.WriteString(scanline)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[ REDACT ] :: error occurred: %v\n", err)
	}

	return results.String()
}

func main() {
	output_unsecure := "Hello, my username is 'hide_this' and my password is 'password123'!!"
	output_expected := "Hello, my username is 'REDACTED' and my password is 'REDACTED'!!"

	fmt.Println("[ Raw .... ] ::", output_unsecure)
	fmt.Println("[ Expected ] ::", output_expected)
	fmt.Println("[ Redacted ] ::", Redact(output_unsecure))
}
