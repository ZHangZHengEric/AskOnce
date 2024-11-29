
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from langchain.document_loaders.excel import UnstructuredExcelLoader

class ExcelLoader(BaseLoader):
    def load_data(self, url):
        loader = UnstructuredExcelLoader(url, mode="elements")
        docs = loader.load()
        return docs
