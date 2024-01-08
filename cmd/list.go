/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	dummy = [][]string{
		{"Hello", "たくみ"},
		{"Godd Morning", "ようこ"},
	}
	outPutFile = "apilist.tsv"
)

// listCmd represents output url and httpmethod list
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Output url and httpmethod list.",
	Long:  "Output url and httpmethod list.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// INFO: 出力用ファイルのオープン
		file, err := os.OpenFile(outPutFile, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return fmt.Errorf("cannot create file: %s", err.Error())
		}
		defer file.Close()

		// INFO: Excelで文字化けしないようにする設定。BOM付きUTF8をfileの先頭に付与
		buf := bufio.NewWriter(file)
		buf.Write([]byte{0xEF, 0xBB, 0xBF})

		// INFO: csv形式でデータを書き込み
		writer := csv.NewWriter(buf)
		writer.Comma = '\t'  //タブ区切りに変更
		defer writer.Flush() //内部バッファのフラッシュは必須
		for _, rec := range dummy {
			if err := writer.Write(rec); err != nil {
				return fmt.Errorf("cannot write record: %s", err.Error())
			}
		}

		fmt.Println("list called")
		return nil
	},
}

func init() {
}
