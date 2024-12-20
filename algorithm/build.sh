#!/bin/bash
export PYTHONPATH="../../:$PYTHONPATH"
# 需要打包的服务文件名，标准的服务是放在相同名称文件夹下的，如：document_split/document_split.py
services=(
 document_split
 answer_by_documents
 data_convert
 generate_outlines
 result_add_reference
 split_question
 question_rewrite
 net_rag_assessment
 search_result_post_process
 )
mkdir -p build
# Loop through the array and execute pyinstaller for each file
for s in "${services[@]}"; do
    echo "服务 $s 构建开始."
    nuitka --nofollow-imports --enable-plugin=multiprocessing --include-package=AskOnce --output-dir=build/  --static-libpython=no -o "build/$s.bin" --remove-output "services/$s/$s.py"
# Wait for all background jobs to complete
    echo "服务 $s 构建完成."
done

