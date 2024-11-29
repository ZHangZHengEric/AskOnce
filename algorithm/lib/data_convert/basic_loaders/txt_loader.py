import os
import hashlib
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader

class TxtLoader(BaseLoader):
    def load_data(self, url: str):
        if not os.path.exists(url):
            raise FileNotFoundError(f"The file at {url} does not exist.")

        with open(url, "r", encoding="utf-8") as file:
            content = file.read()

        id = hashlib.sha256((content + url).encode()).hexdigest()

        meta_data = {"url": url, "file_size": f"{os.path.getsize(url)/ 1024 **2:.02f}MB", "file_type": url.split(".")[-1]}

        return content,{
            "id": id,
            "meta_data": meta_data,
            "data":[
                {
                    "content": content  
                }
            ],
        }
