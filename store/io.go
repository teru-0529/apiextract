package store

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func NewWriter(fileName string) (*csv.Writer, func(), error) {
	dir := filepath.Dir(fileName)

	// INFO: フォルダが存在しない場合作成する
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, nil, fmt.Errorf("cannot create directory: %s", err.Error())
		}
	}

	// INFO: 出力用ファイルのオープン
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create file: %s", err.Error())
	}

	// INFO: Excelで文字化けしないようにする設定。BOM付きUTF8をfileの先頭に付与
	buf := bufio.NewWriter(file)
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	// INFO: tsv形式でデータを書き込み
	writer := csv.NewWriter(buf)
	writer.Comma = '\t' //タブ区切りに変更

	return writer, func() { file.Close() }, nil
}
