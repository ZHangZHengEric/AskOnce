import requests
import json
import logging
logging.basicConfig(level = logging.DEBUG,
                    format='%(asctime)s - %(levelname)s - %(message)s')
logger  = logging.getLogger(__name__)

class SearchService:
    def __init__(self, base_url="http://172.16.0.73:20036",search_engine_api_key=''):
        self.base_url = base_url
        self.headers = {
            "Content-Type": "application/json",
            "User-Source": "AskOnce_bakend"
        }
    # docs [{'title':,'text':,'metadata':{}}]
    def add_data(self, docs: list, collection_name):
        url = f"{self.base_url}/askonce/api/v1/kdb/doc/addByBatchText"
        data = {
            "kdbName": collection_name,
            "autoCreate": True,
            "docs": docs
        }
        response = requests.post(url=url, json=data, headers=self.headers)
        if response.status_code == 200:
            logger.info(f'{collection_name},Documents import successfully!')
            return "Documents import successfully!",json.loads(response.text)
        else:
            logger.error(f"{collection_name},Failed to import documents. Status code: {response.status_code}")
            return f"Failed to import documents. Status code: {response.status_code}", json.loads(response.text)

    def delete_data(self, docNames: list, collection_name):
        url = f"{self.base_url}/open/kdb/col/delete"
        data = {
            "collectionName": collection_name,
            "docNames":[""]
        }

        response = requests.post(url, json=data, headers=self.headers)
        if response.status_code == 200:
            return "Documents deleted successfully!",response.json()
        else:
            return f"Failed to delete documents. Status code: {response.status_code}", response.text

    def delete_collection(self, collection_name: str):
        url = f"{self.base_url}/askonce/api/v1/kdb/delete"
        data = {
            "kdbName": collection_name
        }

        response = requests.post(url, json=data, headers=self.headers)
        if response.status_code == 200:
            logger.info(f'{collection_name}Collection deleted successfully!')
            return "Collection deleted successfully!", response.json()
        else:
            logger.error(f"{collection_name}Failed to delete Collection. Status code: {response.status_code}")
            return f"Failed to delete Collection. Status code: {response.status_code}", response.text

    def search_data(self, text: str, collection_name: str):
        url = f"{self.base_url}/askonce/api/v1/search/kdb"
        data = {
            "kdbName": collection_name,
            "question": text
        }
        logger.debug(f'input data is {data}')
        response = requests.post(url, json=data, headers=self.headers)

        if response.status_code == 200:
            logger.debug('Collection search successfully')
            return 'Collection search successfully',response.text
        else:
            logger.debug(f"Failed to search documents. Status code: {response.status_code}")
            return f"Failed to search documents. Status code: {response.status_code}", response.text


