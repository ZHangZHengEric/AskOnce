import hashlib

try:
    from langchain_community.document_loaders import YoutubeLoader
except ImportError:
    raise ImportError(
        'YouTube video requires extra dependencies. Install with `pip install --upgrade "embedchain[dataloaders]"`'
    ) from None
from AskOnce.algorithm.lib.data_convert.utils.json_serializable import register_deserializable
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from AskOnce.algorithm.lib.data_convert.utils.common_utils import clean_string


@register_deserializable
class YoutubeVideoLoader(BaseLoader):
    def load_data(self, url):
        """Load data from a Youtube video."""
        loader = YoutubeLoader.from_youtube_url(url, add_video_info=True)
        doc = loader.load()
        output = []
        if not len(doc):
            raise ValueError(f"No data found for url: {url}")
        content = doc[0].page_content
        content = clean_string(content)
        metadata = doc[0].metadata
        metadata["url"] = url

        output.append(
            {
                "content": content,
                "meta_data": metadata,
            }
        )
        doc_id = hashlib.sha256((content + url).encode()).hexdigest()
        return {
            "doc_id": doc_id,
            "data": output,
        }
        
        
