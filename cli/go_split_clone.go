package cli

import (
	"fmt"
	"os"
)

const (
	ExitOK int = 0
	ExitNG int = 1
)

func Split(args []string) int {
	cli := &CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintln(cli.Stderr, "Error:", err)
		return ExitNG
	}

	return ExitOK
}

func getSuffix(suffixSlice []rune) []rune {
	lastIndex := len(suffixSlice) - 1
	checkIndex := (len(suffixSlice) / 2) - 1
	if len(suffixSlice) == 0 {
		suffixSlice = append(suffixSlice, rune('a'), rune('a'))
		return suffixSlice
	}
	if string(suffixSlice[lastIndex]) != "z" { // 1文字目が「z」ではない場合
		// 1文字目をインクリメント
		suffixSlice[lastIndex]++
		return suffixSlice
	}

	if string(suffixSlice[checkIndex]) == "y" { // 左からz以外の最初の文字が「y」の場合
		suffixSlice[checkIndex]++
		// 2文字追加
		suffixSlice = append(suffixSlice, rune('a'), rune('a'))
		// 残りを「a」にする
		for index := range suffixSlice {
			if index <= checkIndex {
				continue
			}
			suffixSlice[index] = rune('a')
		}
	} else {
		// 文字列の右から「z」を「a」に、「z」ではない最初の文字を次のアルファベットに
		for index := len(suffixSlice) - 1; index >= 0; index-- {
			if string(suffixSlice[index]) == "z" {
				suffixSlice[index] = rune('a') //「a」にする
			} else {
				suffixSlice[index]++ // 次のアルファベットに
				break
			}
		}
	}
	return suffixSlice
}
