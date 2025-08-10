package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ncruces/zenity"
	"golang.org/x/sys/windows/registry"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <target file or directory>\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		// コマンドライン引数が不正な場合はダイアログでエラー表示
		zenity.Error("使用法: zipwithpwd <対象ファイルまたはディレクトリ>")
		os.Exit(1)
	}
	target := flag.Arg(0)
	// fmt.Printf("[info] target: %s\n", target) // コンソール出力を無効化

	zipName := filepath.Base(target) + ".zip"

	password, err := suggestPassword(target)
	if err != nil {
		// エラーダイアログで表示
		zenity.Error("パスワードサジェスト失敗: " + err.Error())
		os.Exit(2)
	}

	// パスワード編集ダイアログを表示
	password, err = getPasswordFromInputBox(password)
	if err != nil {
		// エラーダイアログで表示
		zenity.Error("パスワード入力ダイアログ失敗: " + err.Error())
		os.Exit(2)
	}

	err = createPasswordZip(zipName, target, password)
	if err != nil {
		// エラーダイアログで表示
		zenity.Error("zip作成失敗: " + err.Error())
		os.Exit(2)
	}
	// 成功メッセージをダイアログで表示
	zenity.Info("zip作成成功: " + zipName)
}

// パスワードテンプレート構造体

// zenityのEntryダイアログでパスワードを編集・取得する
func getPasswordFromInputBox(suggested string) (string, error) {
	pw, err := zenity.Entry(
		"推奨パスワードを編集できます。",
		zenity.Title("パスワード入力"),
		zenity.EntryText(suggested),
	)
	if err != nil {
		// キャンセルやエラー時はサジェスト値を使う
		return suggested, nil
	}
	if pw == "" {
		return suggested, nil
	}
	return pw, nil
}

type PasswordTemplate struct {
	Default  string `json:"default"`
	Patterns []struct {
		Pattern  string `json:"pattern"`
		Template string `json:"template"`
	} `json:"patterns"`
}

// パスワードサジェスト関数
func suggestPassword(target string) (string, error) {
	// 1. ユーザプロファイルディレクトリの zipwithpwd.json を優先的に確認
	homeDir, err := os.UserHomeDir()
	var templatePath string
	if err == nil {
		templatePath = filepath.Join(homeDir, "zipwithpwd.json")
		if _, err := os.Stat(templatePath); err == nil {
			// ユーザプロファイルに設定ファイルが存在する場合はそれを使用
			tpl, err := loadPasswordTemplate(templatePath)
			if err == nil {
				base := filepath.Base(target)
				for _, p := range tpl.Patterns {
					matched, _ := regexp.MatchString(p.Pattern, base)
					if matched {
						return fillTemplate(p.Template, base), nil
					}
				}
				return fillTemplate(tpl.Default, base), nil
			}
		}
	}

	// 2. 実行ファイルと同じディレクトリにある設定ファイルを確認
	exePath, err := os.Executable()
	if err != nil {
		return defaultPassword(), nil
	}
	exeDir := filepath.Dir(exePath)
	exeName := filepath.Base(exePath)
	// 拡張子を.jsonに変更
	templateName := strings.TrimSuffix(exeName, filepath.Ext(exeName)) + ".json"
	templatePath = filepath.Join(exeDir, templateName)

	tpl, err := loadPasswordTemplate(templatePath)
	if err != nil {
		// テンプレートファイルがなければデフォルト
		return defaultPassword(), nil
	}
	base := filepath.Base(target)
	for _, p := range tpl.Patterns {
		matched, _ := regexp.MatchString(p.Pattern, base)
		if matched {
			return fillTemplate(p.Template, base), nil
		}
	}
	return fillTemplate(tpl.Default, base), nil
}

func loadPasswordTemplate(path string) (*PasswordTemplate, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tpl PasswordTemplate
	dec := json.NewDecoder(f)
	if err := dec.Decode(&tpl); err != nil {
		return nil, err
	}
	return &tpl, nil
}

func fillTemplate(tpl, basename string) string {
	date := time.Now().Format("20060102")
	tpl = strings.ReplaceAll(tpl, "{{date}}", date)
	tpl = strings.ReplaceAll(tpl, "{{basename}}", strings.TrimSuffix(basename, filepath.Ext(basename)))
	return tpl
}

func defaultPassword() string {
	return fillTemplate("{{date}}", "")
}

// パスワード付きzipファイルを作成
// 7-Zipのパスをレジストリから取得
func get7ZipPath() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\\7-Zip`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("7-Zipがインストールされていません: %v", err)
	}
	defer k.Close()
	path, _, err := k.GetStringValue("Path")
	if err != nil {
		return "", fmt.Errorf("7-ZipのPath値が取得できません: %v", err)
	}
	return filepath.Join(path, "7z.exe"), nil
}

// 7z.exeでパスワード付きzipを作成
func createPasswordZip(zipPath, targetPath, password string) error {
	sevenZip, err := get7ZipPath()
	if err != nil {
		return err
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if info.IsDir() {
		cmd = exec.Command(sevenZip, "a", "-tzip", "-p"+password, "-mem=ZipCrypto", zipPath, targetPath+"\\*")
	} else {
		cmd = exec.Command(sevenZip, "a", "-tzip", "-p"+password, "-mem=ZipCrypto", zipPath, targetPath)
	}
	// コンソール出力を無効化
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}
