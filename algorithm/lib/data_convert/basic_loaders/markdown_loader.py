from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from typing import Optional, Any
import markdown
import hashlib
import os


class MarkdownLoader(BaseLoader):
    
    def load_data(self, file_path):
        """Load markdown document from loader"""
        # meta_data = {"url": file_path, "file_size": f"{os.path.getsize(file_path)/ 1024 **2:.02f}MB", "file_type": file_path.split(".")[-1]}

        # if need_structure == False:
        #     with open(file_path, "r", encoding="utf-8") as f:
        #         markdown_text = f.read()
        # else:
        with open(file_path, 'r', encoding='utf-8') as f:
            markdown_content = f.read()
            markdown_text = markdown.markdown(markdown_content, extensions=['extra', 'smarty', 'nl2br'])
                # assert markdown_content == markdown_text
        # doc_id = hashlib.sha256((markdown_text + file_path).encode()).hexdigest()
        return markdown_text,{}

    

       