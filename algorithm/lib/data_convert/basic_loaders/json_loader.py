
import hashlib
import json
import os
import re
from typing import Union
from  AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
import requests
from AskOnce.algorithm.lib.data_convert.utils.common_utils import is_valid_json_string, clean_string




class JSONLoader(BaseLoader):
    def load_data(self, file_path):
        with open(file_path, "r", encoding="utf-8") as jf:
            try:
                text_detail = json.load(jf)
                text = json.dumps(text_detail,ensure_ascii=False)
            except:
                'json格式无法解析','json格式无法解析'
        return text,text_detail

    
    
    