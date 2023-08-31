package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	var b = flag.String("b", "", "分割するファイルサイズ")
	var l = flag.Int("l", 0, "分割する行数")
	var n = flag.Int("n", 0, "分割するファイルの個数")

	flag.Parse()

	if *b != "" && *l != 0 { // バイト数と行数を両方指定するとエラー
		println("Error: ", "バイト数と行数が両方指定されています。")
		return
	} else if *b == "" && *l == 0 && *n == 0 { // 何も指定されていない場合
		*l = 1000
	}

	// コマンドライン引数を受け取る
	args := flag.Args()
	fileName := ""
	prefix := ""
	if len(args) == 0 { // splitする対象ファイルが指定されていない
		println("Error: ", "対象のファイルが指定されていません。")
		return
	} else if len(args) > 2 { // コマンドライン引数が多すぎる
		println("Error: ", "指定されている引数が多すぎます。")
		return
	} else if len(args) == 1 {
		fileName = args[0]
	} else if len(args) == 2 {
		fileName = args[0]
		prefix = args[1]
	}

	fp, err := os.Open("./" + fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// ファイル情報取得
	fileinfo, staterr := fp.Stat()

	if staterr != nil {
		fmt.Println(staterr)
		return
	}

	// 接尾辞のスライスを初期化
	var suffixSlice []rune

	if *b != "" {
		targetSize, _ := strconv.Atoi(*b)
		fileCount := int(fileinfo.Size())/targetSize + 1
		// 指定したバイトずつファイルを分割する
		for i := 0; i < fileCount; i++ {
			// 指定バイト分のバッファを用意
			buf := make([]byte, targetSize)

			var start int64 = int64(i * targetSize)
			fp.Seek(start, 0)
			data, err := fp.Read(buf)
			if err == io.EOF {
				continue
			}

			suffixSlice = getSuffix(suffixSlice)
			var suffix string
			for index := range suffixSlice {
				suffix += string(suffixSlice[index])
			}
			err2 := os.WriteFile(prefix+suffix, buf[:data], 0755)

			if err != nil {
				panic(err)
			}
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	} else if *l != 0 {
		scanner := bufio.NewScanner(fp)
		output := ""
		lineCount := 0
		i := 0
		for scanner.Scan() {
			output += scanner.Text() + "\n"
			lineCount += 1
			if lineCount%*l == 0 {

				suffixSlice = getSuffix(suffixSlice)
				var suffix string
				for index := range suffixSlice {
					suffix += string(suffixSlice[index])
				}
				err := os.WriteFile(prefix+suffix, []byte(output), 0755)
				i++
				output = ""
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if output != "" {
			suffixSlice = getSuffix(suffixSlice)
			var suffix string
			for index := range suffixSlice {
				suffix += string(suffixSlice[index])
			}
			err2 := os.WriteFile(prefix+suffix, []byte(output), 0755)
			output = ""
			if err2 != nil {
				log.Fatal(err2)
			}
		}

		// スキャン時のエラーをハンドル
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else if *n != 0 {
		fileCount := *n - 1
		targetSize := int(fileinfo.Size()) / *n
		sizeSurplus := int(fileinfo.Size()) % *n
		// 指定した個数ずつファイルを分割する
		for i := 0; i <= fileCount; i++ {
			bufSize := targetSize
			if i == fileCount { // 最後のチャンクに余ったバイトを加える
				bufSize += sizeSurplus
			}
			// 指定バイト分のバッファを用意
			buf := make([]byte, bufSize)

			var start int64 = int64(i * targetSize)
			fp.Seek(start, 0)
			data, err := fp.Read(buf)
			if err == io.EOF {
				continue
			}

			suffixSlice = getSuffix(suffixSlice)
			var suffix string
			for index := range suffixSlice {
				suffix += string(suffixSlice[index])
			}
			err2 := os.WriteFile(prefix+suffix, buf[:data], 0755)

			if err != nil {
				panic(err)
			}
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}

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
