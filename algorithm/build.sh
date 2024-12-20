#!/bin/bash

# 需要打包的服务文件名，标准的服务是放在相同名称文件夹下的，如：document_split/document_split.py
services=("document_split" "answer_by_documents")

# Loop through the array and execute pyinstaller for each file
for s in "${services[@]}"; do
  echo "开始构建服务 $s..."
  (
    python -m nuitka  --nofollow-imports  --enable-plugin=multiprocessing --include-package=AskOnce --output-dir=build/  --static-libpython=no -o "build/$s.bin" --remove-output services/answer_by_documents/answer_by_documents.py
  ) &
# Wait for all background jobs to complete
wait

echo "所有服务构建完成."
done

