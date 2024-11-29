import hashlib
import docx2txt
try:
    from langchain_community.document_loaders import Docx2txtLoader
    from langchain_community.document_loaders import UnstructuredWordDocumentLoader
except ImportError:
    raise ImportError(
        'Docx file requires extra dependencies. Install with `pip install --upgrade "embedchain[dataloaders]"`'
    ) from None
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader

class DocxLoader(BaseLoader):
    def load_data(self, url, meta_data={}):
        loader = Docx2txtLoader(url)
        output = []
        data = loader.load()
        content = data[0].page_content
        
        if len(content)==0:
            loader = UnstructuredWordDocumentLoader(url)
            print('loader doc')
            output = []
            data = loader.load()
            content = data[0].page_content
        
        metadata = data[0].metadata
        output.append({"content": content, "meta_data": metadata})
        id = hashlib.sha256((content + url).encode()).hexdigest()
        return content,{
            "id": id,
            "meta_data":meta_data,
            "data": output,
        }