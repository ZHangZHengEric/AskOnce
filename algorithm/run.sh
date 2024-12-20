#!/bin/bash
#用于跨平台编译
docker build --platform linux/amd64 -t askonce-algorithm-builder .
#编译可执行的包
docker run --rm -v $(dirname $(pwd)):/src -v $(dirname $(pwd))/bin:/src/bin  --platform linux/amd64 askonce-algorithm-builder bash -c "pip install pyinstaller && cd /src/algorithm && ./build.sh"
