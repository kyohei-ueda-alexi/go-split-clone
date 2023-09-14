package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	cases := map[string]struct {
		args    string
		in      string
		wantErr bool
	}{
		"tooManyOptions":    {"main -b 10 -l 2 file.txt", "", hasErr},
		"fileDoesNotExist":  {"main -l 2 aa", "", hasErr},
		"fileIsNotAssigned": {"main -l 2", "", hasErr},
		"tooManyArgs":       {"main -l 2 ../file.txt test. aa", "", hasErr},
		"byteOption":        {"main -b 5 ../file.txt test.", "", noErr},
		"lineOption":        {"main -l 2 ../file.txt test.", "", noErr},
		"chunkOption":       {"main -n 5 ../file.txt test.", "", noErr},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got bytes.Buffer
			cli := &CLI{
				Stdout: &got,
				Stderr: &got,
				Stdin:  strings.NewReader(tt.in),
			}

			args := strings.Split(tt.args, " ")
			err := cli.Run(args)

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			}
		})
	}

	// [Todo] 各ファイル分割用メソッドのテストを追加する
	// [Todo] テスト追加後、各メソッドをリファクタリングしたい。。。

	// [Todo] テストファイルの作成、テスト時に吐き出されるファイルの削除を自動化したい
	// err := os.RemoveAll("test")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
