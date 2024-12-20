#!/bin/bash

# 需要打包的服务文件名，标准的服务是放在相同名称文件夹下的，如：document_split/document_split.py
services=("document_split" "answer_by_documents")
mkdir -p build
# Loop through the array and execute pyinstaller for each file
for s in "${services[@]}"; do
    python -m nuitka  --nofollow-imports  --enable-plugin=multiprocessing --include-package=AskOnce --output-dir=build/  --static-libpython=no -o "build/$s.bin" --remove-output services/answer_by_documents/answer_by_documents.py
# Wait for all background jobs to complete
    echo "服务 $s 构建完成."
done

