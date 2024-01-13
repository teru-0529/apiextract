/*
Copyright © 2024 Teruaki Sato <andrea.pirlo.0529@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/teru-0529/apiextract/model"
	"github.com/teru-0529/apiextract/store"
)

var (
	dummy = [][]string{
		{"Hello", "たくみ"},
		{"Godd Morning", "ようこ"},
	}
	inputFile = "./openapi/orders/openapi.yaml"
)

var outputFile string

// listCmd represents output url and httpmethod list
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Output url and httpmethod list.",
	Long:  "Output url and httpmethod list.",
	RunE: func(cmd *cobra.Command, args []string) error {

		openapi, err := model.NewOpenApi(inputFile)
		if err != nil {
			return err
		}

		fmt.Println(openapi.Openapi())
		// fmt.Println(openapi.Info())
		// fmt.Printf("%#v\n", openapi.Servers())
		// fmt.Printf("%#v\n", openapi.Tags())
		// for _, v := range openapi.Paths() {
		// 	fmt.Printf("%#v\n", v)
		// }

		// INFO: Writerの取得
		writer, cleanup, err := store.NewWriter(outputFile)
		if err != nil {
			return err
		}
		defer cleanup()

		// INFO: 書き込み
		defer writer.Flush() //内部バッファのフラッシュは必須
		writer.Write([]string{
			"tags",
			"path",
			"method",
			"operationId",
			"summary",
			"description",
			"numOfParameter",
			"requestBody",
			"response",
			"hasExternalDocs",
		})
		for _, path := range openapi.Paths() {
			if err := writer.Write(path.ToArray()); err != nil {
				return fmt.Errorf("cannot write record: %s", err.Error())
			}
		}

		fmt.Println("write apilist")
		return nil
	},
}

func init() {
	// INFO:フラグ値を変数にBind
	listCmd.Flags().StringVarP(&outputFile, "out", "O", "apilist.tsv", "output file path")
}
