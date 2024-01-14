# apiextract

`openapi.yaml`に記載された各種情報を抽出するコマンドラインツール

## Install

Githubのリリースサイトから、各環境用の圧縮ファイルをDownloadし任意のフォルダに展開.<br>
コマンドラインから利用可能.

[最新版のリリースページはコチラ](https://github.com/teru-0529/apiextract/releases/latest)

## Usage

`apiextract.exe`を配置したディレクトリ上で、コンソールから

```console
apiextract [command]
```

の形で実行する.

```console
Available Commands:
  help        Help about any command.
  path        Output url and http-method list.
  version     Show semantic version and release date.

Flags:
  -h, --help            help for apiextract
```

### ■ [path]sub command

```console
Usage:
  apiextract path [flags]

Flags:
  -h, --help         help for path
  -I, --in string    input file path (default "./openapi.yaml")
  -O, --out string   output file path (default "./pathlist.tsv")
```

* `openApi.yaml`ファイルから、`パス`、`HTTPメソッド`ごとの情報をリスト化してtsvファイルに出力する.
* 入出力のファイルパスは、フラグで指定可能.
* 抽出される情報は以下
  |no|        項目        |                   概要説明                   |
  |--|------------------|------------------------------------------|
  |1 |     TagList      |             付与したタグ（複数可）のリスト.             |
  |2 |       Path       |        パス文字列.パスパラメーターは`{}`で囲んである.        |
  |3 |      Method      |               HTTP-METHOD.               |
  |4 |   OperationId    |    ユニークなID.Java実装時はインターフェースのメソッド名になる.    |
  |5 |     Summary      |                 APIサマリ.                  |
  |6 |   Description    |                APIの詳細説明.                 |
  |7 |  NumOfParameter  |           Path/Queryなどの入力パラメータ数.           |
  |8 |  HasRequestbody  |               リクエストボディの有無.               |
  |9 |ResponseStatusList|レスポンス時のHTTP-STATUSリスト.この単位でレスポンスの形が定義される. |
  |10|   ExternalDoc    |             別URLの処理詳細ページ有無.              |

## Conditions

* 分割して記載した`openApi.yaml`には未対応.<br>
  `redocly-cli`など別のツールでbundle後に利用.

## Feature

v1.0では`path`コマンドのみの実装済. 以下機能などを随時拡張予定.

* パラメータ/リクエストボディ、レスポンス内容等、詳細な情報のリスト化
* コンポーネント（共通化した項目）の情報をリスト化
* configファイルの利用(フラグとしてコマンド上に記載する情報を事前にconfigファイルへ登録して利用)
