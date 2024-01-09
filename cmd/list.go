/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/teru-0529/apiextract/store"
)

var (
	dummy = [][]string{
		{"Hello", "たくみ"},
		{"Godd Morning", "ようこ"},
	}
	outPutFile string
)

// listCmd represents output url and httpmethod list
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Output url and httpmethod list.",
	Long:  "Output url and httpmethod list.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// INFO: Writerの取得
		writer, cleanup, err := store.NewWriter(outPutFile)
		if err != nil {
			return err
		}
		defer cleanup()

		// INFO: 書き込み
		defer writer.Flush() //内部バッファのフラッシュは必須
		for _, rec := range dummy {
			if err := writer.Write(rec); err != nil {
				return fmt.Errorf("cannot write record: %s", err.Error())
			}
		}

		fmt.Println("write apilist")
		return nil
	},
}

func init() {
	// INFO:フラグ値を変数にBind
	listCmd.Flags().StringVarP(&outPutFile, "out", "O", "apilist.tsv", "output file path")
}
