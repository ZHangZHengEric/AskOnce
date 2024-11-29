from AskOnce.algorithm.lib.data_convert.utils.json_serializable import JSONSerializable

class BaseLoader(JSONSerializable):
    def __init__(self, factory):
        self.factory = factory

    def load_data(self, url):
        """
        Implemented by child classes
        """
        pass