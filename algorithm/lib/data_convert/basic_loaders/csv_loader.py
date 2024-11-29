import os
import csv
import hashlib
import requests
from urllib.parse import urlparse
from AskOnce.algorithm.lib.data_convert.utils.json_serializable import JSONSerializable
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from io import StringIO

class CsvLoader(BaseLoader):

    @staticmethod
    def _detect_delimiter(first_line):
        delimiters = [",", "\t", ";", "|"]
        counts = {delimiter: first_line.count(delimiter) for delimiter in delimiters}
        return max(counts, key=counts.get)

    @staticmethod
    def _get_file_content(file_path):
        # 加个url检测
        url = urlparse(file_path)
        if all([url.scheme, url.netloc]) and url.scheme not in ["file", "http", "https"]:
            raise ValueError("Not a valid URL.")
        if url.scheme in ["http", "https"]:
            response = requests.get(file_path)
            response.raise_for_status()
            return StringIO(response.text)
        elif url.scheme == "file":
            path = url.path
            return open(path, newline="")  # Open the file using the path from the URI
        else:
            return open(file_path, newline="")  # Treat content as a regular file path

    def load_data(self, file_path):
        text_detail = []
        text = ""
        with CsvLoader._get_file_content(file_path) as file:
            # meta_data = {"url": file_path, "file_size": f"{os.path.getsize(file_path)/ 1024 **2:.02f}MB", "file_type": file_path.split(".")[-1]}
            first_line = file.readline()
            delimiter = CsvLoader._detect_delimiter(first_line)
            file.seek(0)  # Reset the file pointer to the start
            reader = csv.DictReader(file, delimiter=delimiter)
            for i, row in enumerate(reader):
                text_detail.append({"content":row})
                row_str = ','.join([f"{field}: {value}" for field, value in row.items()])
                text += row_str
        # doc_id = hashlib.sha256((file_path + " ".join(lines)).encode()).hexdigest()
        return text,text_detail
    