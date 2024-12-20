#!/bin/bash
#用于跨平台编译
docker run --rm -v ./build.sh:/src -v bin:/src/bin  python:3.8 bash -c "pip install pyinstaller && pyinstaller /src/your_script.py"
