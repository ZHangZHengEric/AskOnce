

# 用来兜底的一个loader, 目前需要下载些依赖，第一次运行需要proxychains
import hashlib
from AskOnce.algorithm.lib.data_convert.utils.json_serializable import register_deserializable
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from AskOnce.algorithm.lib.data_convert.utils.common_utils import clean_string


@register_deserializable
class UnstructuredLoader(BaseLoader):
    def load_data(self, url):
        """Load data from an Unstructured file."""
        try:
            from langchain_community.document_loaders import \
                UnstructuredFileLoader
        except ImportError:
            raise ImportError(
                'Unstructured file requires extra dependencies. Install with `pip install --upgrade "embedchain[dataloaders]"`'  # noqa: E501
            ) from None

        loader = UnstructuredFileLoader(url)
        data = []
        all_content = []
        pages = loader.load_and_split()
        if not len(pages):
            raise ValueError("No data found")
        for page in pages:
            content = page.page_content
            content = clean_string(content)
            metadata = page.metadata
            metadata["url"] = url
            data.append(
                {
                    "content": content,
                    "meta_data": metadata,
                }
            )
            all_content.append(content)
        doc_id = hashlib.sha256((" ".join(all_content) + url).encode()).hexdigest()
        return content, {
            "id": doc_id,
            "data": data,
        }