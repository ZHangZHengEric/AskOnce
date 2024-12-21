from AskOnce.algorithm.lib.llm_api.llm_base_api import LLMBaseAPI
import re
import logging
logging.basicConfig(level = logging.DEBUG,
                    format='%(asctime)s - %(levelname)s - %(message)s')
logger  = logging.getLogger(__name__)

class Translate(LLMBaseAPI):
    def translate_one(self,input_text,target_language='中文'):
        prompt='''将下面这段文字翻译成{target_language}
{input_text}
'''
        trans_result = self.ask_llm(prompt.format(target_language=target_language,input_text=input_text),temperature=0.2)
        if trans_result is not None:
            print('翻译结果',trans_result)
            if '我无法回答这个问题' in trans_result or '我无法理解您的问题' in trans_result or '更多上下文' in trans_result or '更多的上下文' in trans_result or '我无法理解您提供的文字' in trans_result:
                return input_text
            return trans_result
        else:
            print('无法翻译',input_text)
            return input_text
        
    def split_chunks(self,input_text,max_input_length):
        parts = re.split(r'([.。\n\r])',input_text)
        chunks = ['']
        for part_one in parts:
            if len(part_one)+len(chunks[-1])>max_input_length:
                if len(part_one)<2:
                    chunks[-1]+=part_one
                else:
                    chunks.append(part_one)
            else:
                chunks[-1]+=part_one
        return chunks
    
    def translate(self,input_text,target_language):
        target_text  = ''
        if len(input_text)<200:
            target_text = self.translate_one(input_text,target_language)
        else:
            chunks = self.split_chunks(input_text,200)            
            for index,one_chunk in enumerate(chunks):
                print('翻译进度',index,'/',len(chunks))
                target_text += (self.translate_one(one_chunk) + ('\n' if one_chunk.endswith('\n') else "" ))
        return target_text