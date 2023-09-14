package cli

import (
	"testing"
)

// テーブル駆動テストの手法に沿ったコード
func TestGetSuffix(t *testing.T) {
	var tests = []struct {
		input    []rune
		expected string
	}{
		{[]rune{}, "aa"},
		{[]rune{rune('a'), rune('a')}, "ab"},
		{[]rune{rune('y'), rune('z')}, "zaaa"},
		{[]rune{rune('z'), rune('a'), rune('a'), rune('a')}, "zaab"},
		{[]rune{rune('z'), rune('a'), rune('a'), rune('z')}, "zaba"},
	}

	for _, test := range tests {
		output := getSuffix(test.input)
		var outputStr string
		for index := range output {
			outputStr += string(output[index])
		}
		if outputStr != test.expected {
			t.Error("Test Failed: ", "input/", test.input, " expected/", test.expected, " output/", outputStr)
		}
	}
}
