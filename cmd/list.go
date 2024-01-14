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

var inputFile string
var outputFile string

// listCmd represents output url and httpmethod list
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Output url and httpmethod list.",
	Long:  "Output url and httpmethod list.",
	RunE: func(cmd *cobra.Command, args []string) error {

		_, paths, _, err := model.NewOpenApi(inputFile)
		if err != nil {
			return err
		}

		// INFO: Writerの取得
		writer, cleanup, err := store.NewWriter(outputFile)
		if err != nil {
			return err
		}
		defer cleanup()

		// INFO: 書き込み
		defer writer.Flush() //内部バッファのフラッシュは必須
		writer.Write([]string{
			"tagList",
			"path",
			"method",
			"operationId",
			"summary",
			"description",
			"numOfParameter",
			"hasRequestbody",
			"responseStatusList",
			"hasExternalDoc",
		})
		for _, path := range *paths {
			if err := writer.Write(path.ToPathArray()); err != nil {
				return fmt.Errorf("cannot write record: %s", err.Error())
			}
		}

		fmt.Println("write path list")
		return nil
	},
}

func init() {
	// INFO:フラグ値を変数にBind
	listCmd.Flags().StringVarP(&inputFile, "in", "I", "./openapi.yaml", "input file path")
	listCmd.Flags().StringVarP(&outputFile, "out", "O", "./pathlist.tsv", "output file path")
}
