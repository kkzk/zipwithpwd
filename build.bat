@echo off
echo ====================================
echo zipwithpwd ビルド＆パッケージング
echo ====================================

REM 1. buildディレクトリを作成（存在しない場合）
if not exist "build" mkdir build

REM 2. Goでビルド（出力先をbuildフォルダに指定）
echo [1/3] Goビルド中...
go build -o build\zipwithpwd.exe
if errorlevel 1 (
    echo エラー: Goビルドに失敗しました
    pause
    exit /b 1
)
echo Goビルド完了: build\zipwithpwd.exe

REM 3. 設定ファイルもbuildフォルダにコピー（存在する場合）
if exist "zipwithpwd.json" (
    echo [2/3] 設定ファイルをコピー中...
    copy zipwithpwd.json build\zipwithpwd.json >nul
    echo 設定ファイルコピー完了
) else (
    echo [2/3] 設定ファイルが見つかりません（スキップ）
)

REM 4. NSISインストーラー作成
echo [3/3] インストーラー作成中...
makensis installer.nsi
if errorlevel 1 (
    echo エラー: インストーラー作成に失敗しました
    echo NSISがインストールされているか確認してください
    pause
    exit /b 1
)
echo インストーラー作成完了: zipwithpwd_installer.exe

echo ====================================
echo 完了しました！
echo - 実行ファイル: build\zipwithpwd.exe
echo - インストーラー: zipwithpwd_installer.exe
echo ====================================
pause
