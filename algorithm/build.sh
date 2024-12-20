#!/bin/bash

# 需要打包的服务文件名，标准的服务是放在相同名称文件夹下的，如：document_split/document_split.py
files=("document_split" "answer_by_documents")

# Loop through the array and execute pyinstaller for each file
for file in "${files[@]}"; do
  echo "开始构建服务 $file..."
  (
    pyinstaller --onefile --strip --noconfirm --distpath ./bin/ --add-data "../algorithm:AskOnce/algorithm" --name "$file" "services/$file/$file.py"
  ) &
# Wait for all background jobs to complete
wait

echo "所有服务构建完成."
done
