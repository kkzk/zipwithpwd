@echo off
echo ====================================
echo zipwithpwd �r���h���p�b�P�[�W���O
echo ====================================

REM 1. build�f�B���N�g����dist�f�B���N�g�����쐬�i���݂��Ȃ��ꍇ�j
if not exist "build" mkdir build
if not exist "dist" mkdir dist

REM 2. Go�Ńr���h�i�o�͐��build�t�H���_�Ɏw��j
echo [1/3] Go�r���h��...
go build -ldflags="-H windowsgui" -o build\zipwithpwd.exe
if errorlevel 1 (
    echo �G���[: Go�r���h�Ɏ��s���܂���
    pause
    exit /b 1
)
echo Go�r���h����: build\zipwithpwd.exe

REM 3. �ݒ�t�@�C����build�t�H���_�ɃR�s�[�i���݂���ꍇ�j
if exist "zipwithpwd.json" (
    echo [2/3] �ݒ�t�@�C�����R�s�[��...
    copy zipwithpwd.json build\zipwithpwd.json >nul
    echo �ݒ�t�@�C���R�s�[����
) else (
    echo [2/3] �ݒ�t�@�C����������܂���i�X�L�b�v�j
)

REM 4. NSIS�C���X�g�[���[�쐬
echo [3/3] �C���X�g�[���[�쐬��...
makensis installer.nsi
if errorlevel 1 (
    echo �G���[: �C���X�g�[���[�쐬�Ɏ��s���܂���
    echo NSIS���C���X�g�[������Ă��邩�m�F���Ă�������
    pause
    exit /b 1
)
echo �C���X�g�[���[�쐬����: zipwithpwd_installer.exe

echo ====================================
echo �������܂����I
echo - ���s�t�@�C��: build\zipwithpwd.exe
echo - �C���X�g�[���[: dist\zipwithpwd_installer.exe
echo ====================================

