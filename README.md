# zipwithpwd

パスワード付きZIPファイルを簡単に作成するWindowsアプリケーションです。

## 概要

zipwithpwdは、ファイルやディレクトリを7-Zipを使用してパスワード付きZIPファイルに圧縮するGUI付きのコマンドラインツールです。日付やファイル名を基にしたパスワードテンプレート機能により、一貫性のあるパスワード管理を支援します。

## 機能

- ファイル・ディレクトリのパスワード付きZIP圧縮
- GUI形式のパスワード入力ダイアログ
- カスタマイズ可能なパスワードテンプレート
- ファイル拡張子に応じた自動パスワード生成
- 7-Zipエンジンによる高品質な圧縮

## 必要な環境

- Windows OS
- [7-Zip](https://www.7-zip.org/) がインストールされていること

## インストール

1. リリースページから最新の実行ファイル (`zipwithpwd.exe`) をダウンロード
2. 任意のフォルダに配置
3. パスワードテンプレートファイル (`zipwithpwd.json`) を設定（オプション）

## 使用方法

### 基本的な使い方

```powershell
zipwithpwd <対象ファイルまたはディレクトリ>
```

例:
```powershell
# ファイルを圧縮
zipwithpwd document.xlsx

# ディレクトリを圧縮
zipwithpwd C:\Users\username\Documents\project
```

### 動作の流れ

1. コマンドラインで対象ファイル/ディレクトリを指定
2. パスワードテンプレートに基づいて推奨パスワードが生成される
3. パスワード入力ダイアログが表示され、推奨パスワードを編集可能
4. 7-Zipを使用してパスワード付きZIPファイルを作成
5. 成功/失敗のダイアログメッセージが表示される

## パスワードテンプレートの設定

設定ファイル `zipwithpwd.json` を以下の場所に配置できます：

1. **ユーザープロファイル** (`%USERPROFILE%\zipwithpwd.json`) - 優先度高
2. **実行ファイルと同じディレクトリ** (`zipwithpwd.json`) - 優先度低

### 設定ファイルの例

```json
{
  "default": "zipwithpwd{{date}}",
  "patterns": [
    {
      "pattern": ".*\\.xlsx$",
      "template": "zipwithpwd-{{basename}}-{{date}}"
    },
    {
      "pattern": ".*\\.docx$",
      "template": "zipwithpwd-{{basename}}-{{date}}"
    }
  ]
}
```

### テンプレート変数

- `{{date}}` - 現在の日付（YYYYMMDD形式）
- `{{basename}}` - 対象ファイル/ディレクトリ名（拡張子なし）

### 設定の詳細

- **default**: デフォルトのパスワードテンプレート
- **patterns**: ファイル拡張子ごとの個別設定
  - **pattern**: 正規表現でファイル名をマッチング
  - **template**: マッチした場合に使用するパスワードテンプレート

## ビルド方法

Go 1.24.4以降が必要です。

```powershell
# 依存関係のインストール
go mod download

# ビルド
go build -o zipwithpwd.exe

# または build.bat を実行
build.bat
```

## 依存ライブラリ

- [github.com/ncruces/zenity](https://github.com/ncruces/zenity) - GUIダイアログ
- [golang.org/x/sys/windows/registry](https://pkg.go.dev/golang.org/x/sys/windows/registry) - Windowsレジストリアクセス

## トラブルシューティング

### よくある問題

**「7-Zipがインストールされていません」エラー**
- 7-Zipをインストールしてください

**「パスワード入力ダイアログ失敗」エラー**
- zenityの初期化に失敗している可能性があります
- アプリケーションを管理者権限で実行してみてください

**ZIPファイルが作成されない**
- 対象ファイル/ディレクトリが存在することを確認してください
- 書き込み権限があることを確認してください
- ファイル名に特殊文字が含まれていないか確認してください

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 作成者

開発者: kkzk
