
import sys
import logging
import hashlib
import requests

from bs4 import BeautifulSoup
from AskOnce.algorithm.lib.data_convert.utils.common_utils import clean_string, is_valid_json_string
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from AskOnce.algorithm.lib.data_convert.utils.json_serializable import register_deserializable

@register_deserializable
class URLLoader(BaseLoader):
    # Shared session for all instances
    _session = requests.Session()

    def load_data(self, url):
        """Load data from a web page using a shared requests' session."""
        headers = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",  # noqa:E501
        }
        response = self._session.get(url, headers=headers, timeout=30)
        response.raise_for_status()
        data = response.content
        content = self._get_clean_content(data, url)

        metadata = {"url": url}

        doc_id = hashlib.sha256((content + url).encode()).hexdigest()
        return {
            "doc_id": doc_id,
            "data": [
                {
                    "content": content,
                    "meta_data": metadata,
                }
            ],
        }

    @staticmethod
    def _get_clean_content(html, url) -> str:
        soup = BeautifulSoup(html, "html.parser")
        original_size = len(str(soup.get_text()))

        tags_to_exclude = [
            "nav",
            "aside",
            "form",
            "header",
            "noscript",
            "svg",
            "canvas",
            "footer",
            "script",
            "style",
        ]
        for tag in soup(tags_to_exclude):
            tag.decompose()

        ids_to_exclude = ["sidebar", "main-navigation", "menu-main-menu"]
        for id_ in ids_to_exclude:
            tags = soup.find_all(id=id_)
            for tag in tags:
                tag.decompose()

        classes_to_exclude = [
            "elementor-location-header",
            "navbar-header",
            "nav",
            "header-sidebar-wrapper",
            "blog-sidebar-wrapper",
            "related-posts",
        ]
        for class_name in classes_to_exclude:
            tags = soup.find_all(class_=class_name)
            for tag in tags:
                tag.decompose()

        content = soup.get_text()
        content = clean_string(content)

        cleaned_size = len(content)
        if original_size != 0:
            logging.info(
                f"[{url}] Cleaned page size: {cleaned_size} characters, down from {original_size} (shrunk: {original_size-cleaned_size} chars, {round((1-(cleaned_size/original_size)) * 100, 2)}%)"  # noqa:E501
            )

        return content

    @classmethod
    def close_session(cls):
        cls._session.close()