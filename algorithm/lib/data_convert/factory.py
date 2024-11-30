
import os
# import magic
import json
import re
from typing import Optional
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.txt_loader import TxtLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.pdf_loader import PdfLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.csv_loader import CsvLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.email_loader import EmailLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.docx_loader import DocxLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.doc_loader import DocLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.markdown_loader import MarkdownLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.pptx_loader import PptxLoader
from AskOnce.algorithm.lib.data_convert.basic_loaders.json_loader import JSONLoader
class BasicLoaderType:
    TXT = TxtLoader
    PDF = PdfLoader
    DOCX = DocxLoader
    DOC = DocLoader
    EMAIL = EmailLoader
    CSV = CsvLoader
    MARKDOWN = MarkdownLoader
    JSON = JSONLoader
    

def guess_file_type(file_path):
    file_extension = file_path.split(".")[-1]
    print(file_extension) 
    if file_extension in ['docx']:
        return "DOCX"
    elif file_extension in ['doc']:
        return "DOC"
    elif file_extension in ['pdf']:
        return "PDF"
    elif file_extension in ['txt']:
        return "TXT"
    elif file_extension in ['csv']:
        return "CSV"
    elif file_extension in ['eml']:
        return "EMAIL"
    elif file_extension in ['md','MD']:
        return "MARKDOWN"
    elif file_extension in ['json']:
        return "JSON"


class BasicLoaderFactory(BaseLoader):
    def __init__(self,correct_dict_path):
        self.correct_dict_path = correct_dict_path
        if os.path.exists(self.correct_dict_path):
            self.correct_dict = json.load(open(correct_dict_path))
            print("correct_dict资源加载suss")
        else:
            self.correct_dict = {}
            
    def create(self, file_path, is_remove_wrap=False):
        
        file_type = guess_file_type(file_path)
        # print('file_type',file_type)
        if file_type:
            basic_loader_type = getattr(BasicLoaderType, file_type)
            text,text_detail = basic_loader_type(self).load_data(file_path)
            text = self.replace_wrong_char(text)
            text = self.remove_duplicate_char(text,is_remove_wrap)
            return text,text_detail
        else:
            print("当前文件类型不支持 loader")
            # raise Exception(f"{basic_loader_type}数据加载器初始化出错")
    
    def replace_wrong_char(self,text):
        for index, one_char in enumerate(text):
            if one_char in self.correct_dict:
                # text['index'] = self.correct_dict[one_char]
                char_list = list(text)
                char_list[index] = self.correct_dict[one_char]
                text = "".join(char_list)
        return text
    
    def remove_duplicate_char(self,text,is_remove_wrap=False):
        text = re.sub(r"\n+",'\n', text)
        if is_remove_wrap:
            text = text.replace('\n','。')
            text = text.replace('.。','。')
            text = re.sub(r"。+",'。', text)
        text = re.sub(r" +",' ', text)
        text = re.sub(r"　+",' ', text)
        return text
    
    