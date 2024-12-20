#!/bin/bash
#用于跨平台编译
docker run --rm -v $(dirname $(pwd)):/src -v $(dirname $(pwd))/bin:/src/bin  --platform linux/amd64 python:3.10 bash -c "pip install pyinstaller && cd /src/algorithm && ./build.sh"
