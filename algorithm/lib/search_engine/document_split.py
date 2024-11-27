import re
import hashlib    

class TextChunkPorcessor:
    def __init__(
        self, 
    ):
        pass
    
    def genenrate_id_based_string(self,text, hash_fun='md5'):
        '''
            根据字符串的内容，生成一个唯一id, 内容 -> id 是单射
        '''
        if hash_fun == 'md5':
            hash_object = hashlib.md5(text.encode())
            unique_id = hash_object.hexdigest()
        elif hash_fun == 'sha_256':
            # 更安全点
            hash_object = hashlib.sha256(text.encode())
            unique_id = hash_object.hexdigest()
        return unique_id
    
    def cutting_by_punctuation(self,text, doc_name):
        '''
            长度太长了 bug需要修复
        '''
        pattern = r"[？！。]"
        sentences_list = []
        sentences = re.split(pattern, text)
        for i, s in enumerate(sentences):
            if i == 0:
                start, end = 0, len(s)
            else:
                start = end
                end = start + len(s)
            if len(s) > 0:
                sentences_list.append({"doc_name":doc_name, "text":s, "start":start, "end":end})
        return sentences_list

    # 按照标点符号，切割后的句子，在length周围
    def segment_length_by_punctuation(self,text_: str, length_: int):
        pattern = re.compile(r"[.。]\n?")
        positions = list(pattern.finditer(text_))
        inner_sentences = []
        start = 0
        for pos in positions:
            end = pos.end()
            if end - start > length_ or end == len(text_):
                inner_sentences.append(  {"passage_id": self.genenrate_id_based_string(text_[start:end]), "passage_content": text_[start:end],"start": start,"end":end})
                start = end
        if start < len(text_):
            inner_sentences.append(  {"passage_id": self.genenrate_id_based_string(text_[start:]), "passage_content": text_[start:],"start": start,"end":len(text_)})
        return inner_sentences

    def merge_sentences_split(self,doc_content,fix_length_list=[128,256,512]):
        '''
        切割逻辑是：先按照标点.。切割,将text分成一个个独立句子，
        当test:end-start大于length,就合并。
        '''
        if not isinstance(fix_length_list, List):
            raise Exception("固定长度参数类型错误，必须为整数列表")
        for fl in fix_length_list:
            if not isinstance(fl, int):
                raise Exception("固定长度列表元素参数类型错误，列表中的元素必须为整数")
        
        if fix_length_list is None:
            fix_length_list = [512]

        sentences = []
        for length in fix_length_list:
            sentences += self.segment_length_by_punctuation(doc_content, length)
        return sentences
    
    # 固定长度滑动窗口分割
    def move_window_split(self, doc_content, window_size = 256, stride= 170):
        text_content = doc_content
        
        start, end , sentent_list = 0, window_size, []
        content_len = len(text_content)
        if content_len <= window_size:
            passage_id = self.genenrate_id_based_string(text_content)
            passage_content = text_content
            sentent_list.append(
                {
                    "passage_id": passage_id, 
                    "passage_content": passage_content,
                    "start": 0,
                    "end":content_len
                }
            )
            return sentent_list
        
        while end < content_len: 
            passage_content = text_content[start:end]
            passage_id = self.genenrate_id_based_string(passage_content)
            sentent_list.append(
                {
                    "passage_id": passage_id,
                    "passage_content": passage_content, 
                    "start": start,
                    "end": end
                }
            )
            start += stride
            end = start + window_size
        # 扫尾
        passage_content = text_content[start:]
        passage_id = self.genenrate_id_based_string(passage_content)
        end = content_len
        sentent_list.append(
            {
                "passage_id":passage_id,
                "passage_content":passage_content,
                "start":start, 
                "end": end
            }
        )
        return sentent_list