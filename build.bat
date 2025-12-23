@echo off
chcp 65001 >nul
REM 设置脚本在遇到错误时停止执行
setlocal enabledelayedexpansion

REM 获取当前目录（根目录）
set ROOT_DIR=%cd%

echo ========================================
echo 开始构建项目
echo ========================================


echo 正在构建wails本...
wails build
if errorlevel 1 (
    echo 错误：wails版本构建失败
    goto end
)

echo wails项目构建完成！

REM 3. 复制指定文件和文件夹到build目录
echo 正在复制必要文件到build目录...
if not exist "build/bin" mkdir "build/bin"


REM 复制文件夹
echo 正在复制 ffmpeg 目录...
xcopy /E /I /Y "ffmpeg" "build\bin\ffmpeg" >nul 2>&1
if errorlevel 1 (
    echo 警告：ffmpeg目录复制失败或不存在
)
echo 正在复制 msyh.ttc 字体文件...
xcopy /Y "msyh.ttc" "build\bin" >nul 2>&1
if errorlevel 1 (
    echo 警告：msyh.ttc文件复制失败或不存在
)


echo ========================================
echo 所有构建任务已完成！
echo ========================================

:end
echo.
echo 按任意键关闭窗口...
pause >nul