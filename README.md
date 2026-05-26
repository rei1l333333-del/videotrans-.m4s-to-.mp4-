 # videotrans

  视频转换 & 字幕生成桌面工具，基于 Go + Fyne 开发。

  ## 功能

  - **M4S → MP4 转换** — FFmpeg 转封装，不重新编码，速度极快
  - **AI 字幕生成** — 基于 OpenAI Whisper，支持中/英/日/韩等多语言
  - 暗色主题 GUI，支持拖拽文件

  ## 依赖

  ### 1. Git

  后续步骤需要用到 Git，请先安装。

  **Windows：**

  前往 https://git-scm.com/download/win 下载安装包，安装时一路默认即可。

  安装完成后打开终端验证：

  ```bash
  git --version

  2. Go

  项目使用 Go 编写，需要 Go 1.21+。

  前往 https://go.dev/dl/ 下载安装包。

  验证：

  go version

  3. FFmpeg

  用于 M4S → MP4 转封装。如果未安装，程序会回退到二进制复制模式（不推荐）。

  Windows（推荐 - 使用 winget）：

  winget install Gyan.FFmpeg

  或手动安装：

  1. 前往 https://www.gyan.dev/ffmpeg/builds/ 下载 ffmpeg-release-essentials.zip
  2. 解压到任意目录（如 C:\ffmpeg）
  3. 将 C:\ffmpeg\bin 添加到系统环境变量 PATH

  验证：

  ffmpeg -version

  4. Python & pip

  Whisper 依赖 Python 环境，需要 Python 3.8+。

  前往 https://www.python.org/downloads/ 下载安装包，安装时务必勾选 "Add Python to PATH"。

  验证：

  python --version
  pip --version

  5. OpenAI Whisper

  用于 AI 字幕生成。

  通过 pip 安装（推荐）：

  pip install openai-whisper

  或通过 Git 源码安装（获取最新版）：

  pip install git+https://github.com/openai/whisper.git

  ▎ Whisper 依赖 PyTorch，pip 会自动安装。如果你有 NVIDIA 显卡并希望使用 GPU 加速，请先安装 CUDA 版 PyTorch：
  ▎
  ▎ pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121
  ▎ pip install openai-whisper

  验证：

  whisper --help

  构建

  git clone https://github.com/rei1l333333-del/videotrans-.m4s-to-.mp4-.git
  cd videotrans-.m4s-to-.mp4-
  go build -o videotrans.exe

  使用

  直接运行：

  ./videotrans.exe

  - 视频转换 标签页：选择或拖拽 .m4s 文件，点击"开始转换"
  - 字幕生成 标签页：选择或拖拽 MP4/音频文件，选择模型和语言，点击"生成字幕"

  字幕模型说明

  ┌────────┬──────────┬──────┬──────┐
  │  模型  │   大小   │ 速度 │ 精度 │
  ├────────┼──────────┼──────┼──────┤
  │ tiny   │ ~39 MB   │ 最快 │ 低   │
  ├────────┼──────────┼──────┼──────┤
  │ base   │ ~74 MB   │ 快   │ 中   │
  ├────────┼──────────┼──────┼──────┤
  │ small  │ ~244 MB  │ 中   │ 较高 │
  ├────────┼──────────┼──────┼──────┤
  │ medium │ ~769 MB  │ 慢   │ 高   │
  ├────────┼──────────┼──────┼──────┤
  │ large  │ ~1550 MB │ 最慢 │ 最高 │
  └────────┴──────────┴──────┴──────┘

  ▎ 首次使用某个模型时会自动下载，请耐心等待。

  许可证

  MIT

  涵盖了 Git、Go、FFmpeg、Python/pip、Whisper 的完整安装命令，包括 pip 安装和 git 源码安装两种方式，还加了 GPU
  加速提示和模型对比表。你可以根据需要调整细节。
