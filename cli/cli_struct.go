package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

type InputFile struct {
	file   *os.File
	size   int64
	name   string
	prefix string
}

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func (cli *CLI) Run(args []string) error {
	var (
		b string
		l int
		n int
	)
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagSet.StringVar(&b, "b", "", "分割するファイルサイズ")
	flagSet.IntVar(&l, "l", 0, "分割する行数")
	flagSet.IntVar(&n, "n", 0, "分割するファイルの個数")

	err := flagSet.Parse(args[1:])
	if err != nil {
		return err
	}

	optionCount := 0
	if b != "" {
		optionCount++
	}
	if l != 0 {
		optionCount++
	}
	if n != 0 {
		optionCount++
	}

	if optionCount >= 2 { // オプション指定が2つ以上
		return fmt.Errorf("指定しているオプションが多すぎます")
	} else if b == "" && l == 0 && n == 0 { // 何も指定されていない場合
		l = 1000
	}

	fileName := ""
	prefix := ""
	if flagSet.NArg() == 0 { // splitする対象ファイルが指定されていない
		return fmt.Errorf("対象のファイルが指定されていません。")
	} else if flagSet.NArg() > 2 { // コマンドライン引数が多すぎる
		return fmt.Errorf("指定されている引数が多すぎます。")
	} else if flagSet.NArg() == 1 {
		fileName = flagSet.Args()[0]
	} else if flagSet.NArg() == 2 {
		fileName = flagSet.Args()[0]
		prefix = flagSet.Args()[1]
	}

	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	// ファイル情報取得
	fileinfo, staterr := fp.Stat()

	if staterr != nil {
		return fmt.Errorf("対象のファイルが存在しません。")
	}

	inputFile := &InputFile{
		file:   fp,
		size:   fileinfo.Size(),
		name:   fileName,
		prefix: prefix,
	}

	if b != "" {
		err := splitByByte(&b, inputFile)
		if err != nil {
			return err
		}
	} else if l != 0 {
		err := splitByLine(&l, inputFile)
		if err != nil {
			return err
		}
	} else if n != 0 {
		err := splitByChunk(&n, inputFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitByByte(b *string, inputFile *InputFile) error {
	// 接尾辞のスライスを初期化
	var suffixSlice []rune

	targetSize, _ := strconv.Atoi(*b)
	fileCount := int(inputFile.size)/targetSize + 1
	// 指定したバイトずつファイルを分割する
	for i := 0; i < fileCount; i++ {
		// 指定バイト分のバッファを用意
		buf := make([]byte, targetSize)

		var start int64 = int64(i * targetSize)
		inputFile.file.Seek(start, 0)
		data, err := inputFile.file.Read(buf)
		if err == io.EOF {
			continue
		}
		if err != nil {
			return err
		}

		suffixSlice = getSuffix(suffixSlice)
		var suffix string
		for index := range suffixSlice {
			suffix += string(suffixSlice[index])
		}
		if err := os.WriteFile(inputFile.prefix+suffix, buf[:data], 0755); err != nil {
			return err
		}
	}

	return nil
}

func splitByLine(l *int, inputFile *InputFile) error {
	// 接尾辞のスライスを初期化
	var suffixSlice []rune

	scanner := bufio.NewScanner(inputFile.file)
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
			err := os.WriteFile(inputFile.prefix+suffix, []byte(output), 0755)
			i++
			output = ""
			if err != nil {
				return err
			}
		}
	}

	if output != "" {
		suffixSlice = getSuffix(suffixSlice)
		var suffix string
		for index := range suffixSlice {
			suffix += string(suffixSlice[index])
		}
		err := os.WriteFile(inputFile.prefix+suffix, []byte(output), 0755)
		output = ""
		if err != nil {
			return err
		}
	}

	// スキャン時のエラーをハンドル
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func splitByChunk(n *int, inputFile *InputFile) error {
	// 接尾辞のスライスを初期化
	var suffixSlice []rune

	fileCount := *n - 1
	targetSize := int(inputFile.size) / *n
	sizeSurplus := int(inputFile.size) % *n
	// 指定した個数ずつファイルを分割する
	for i := 0; i <= fileCount; i++ {
		bufSize := targetSize
		if i == fileCount { // 最後のチャンクに余ったバイトを加える
			bufSize += sizeSurplus
		}
		// 指定バイト分のバッファを用意
		buf := make([]byte, bufSize)

		var start int64 = int64(i * targetSize)
		inputFile.file.Seek(start, 0)
		data, err := inputFile.file.Read(buf)
		if err == io.EOF {
			continue
		}
		if err != nil {
			return err
		}

		suffixSlice = getSuffix(suffixSlice)
		var suffix string
		for index := range suffixSlice {
			suffix += string(suffixSlice[index])
		}
		if err := os.WriteFile(inputFile.prefix+suffix, buf[:data], 0755); err != nil {
			return err
		}
	}
	return nil
}
