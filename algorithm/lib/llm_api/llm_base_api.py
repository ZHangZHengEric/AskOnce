import time
from typing import Any

from openai import OpenAI
import requests
import traceback
class LLMBaseAPI:
    def __init__(
            self, 
            platform_api_url: str , 
            api_key: str ,
            model_name: str ,
            search_url: str = None
        ):
        self.platform_api_url = platform_api_url
        self.api_key  =api_key
        self.model_name = model_name
        self.search_url = search_url
        self.client = OpenAI(
            api_key=api_key,
            base_url=platform_api_url
        )    
        
    
    def search_internet(self,queston,search_session_id):
        if self.search_url is not None:
            url = self.search_url+'/askonce/api/v1/search/web'
            headers = {'Content-Type': 'application/json'}
            data = {
                "sessionId":search_session_id,
                "question" : queston,
            }
            try:
                response = requests.post(url, headers=headers, json=data)
                return response.json()["data"]['data']
            except:
                traceback.print_exc()
                print('检查askonce 网络搜索引擎服务')
        else:
            raise ValueError('set the search url first')
    
    
    def ask_llm(
            self, 
            prompt: str | list, 
            temperature: float = 0.01, 
            presence_penalty: float = 0, 
            frequency_penalty: float = 0,
            max_tokens: int = 1024, 
        ):
        if isinstance(prompt, str):
            messages = [
                {"role": "user", "content": prompt}
            ]
        else:
            messages = prompt 
        start_time = time.time()
        completion = self.client.chat.completions.create(
            model=self.model_name,
            messages=messages,
            temperature=temperature,
            presence_penalty=presence_penalty,
            frequency_penalty=frequency_penalty,
            max_tokens=max_tokens,
            stream=False
        ) 
        return_content = completion.choices[0].message.content
        print('使用模型', self.model_name, '推理耗时', f'{time.time() - start_time:.4f} 秒')
        return return_content
    
    def ask_llm_stream(
            self, 
            prompt: str | list, 
            temperature: float = 0.01, 
            presence_penalty: float = 0, 
            frequency_penalty: float = 0,
            max_tokens: int = 1024, 
        ):
        if isinstance(prompt, str):
            messages = [
                {"role": "user", "content": prompt}
            ]
        else:
            messages = prompt 
        start_time = time.time()
        completion = self.client.chat.completions.create(
            model=self.model_name,
            messages=messages,
            temperature=temperature,
            presence_penalty=presence_penalty,
            frequency_penalty=frequency_penalty,
            max_tokens=max_tokens,
            stream=True
        )
        content = ''
        for chunk in completion:
	# 在这里，每个 chunk 的结构都与之前的 completion 相似，但 message 字段被替换成了 delta 字段
            delta = chunk.choices[0].delta
            if delta.content:
                content= delta.content
                yield content
        print('使用模型', self.model_name, '推理耗时', f'{time.time() - start_time:.4f} 秒')
        return content

    
    def select_search_result(self,search_result_origin,min_chars=50,max_chars=7000,max_search_content_length = 10000):
        all_search_result = []
        search_result=[]
        j=0
        all_content_length = 0
        for i in range(len(search_result_origin)):
            if 'htmlText' not in search_result_origin[i].keys():
                search_result_origin[i]['htmlText'] = search_result_origin[i]['content']
            search_content_max = max(len(search_result_origin[i]['content']),len(search_result_origin[i]['htmlText']))
            if search_content_max>min_chars and search_content_max<max_chars:
                search_result.append(search_result_origin[i])
                all_content_length += search_content_max
                j+=1
            if all_content_length>max_search_content_length:
                break
            if j>10:
                break
        for search_result_index in range(len(search_result)):
            keywords = ['为保证您的正常访问，请输入验证码进行验证。开始验证', '阻止了您进一步访问', '腾讯云EdgeOne','网络安全策略']
            if any(keyword in search_result[search_result_index]['content'] for keyword in keywords):
                continue
            search_content = ''
            search_content = search_result[search_result_index]['content'] if len(search_result[search_result_index]['htmlText']) <len(search_result[search_result_index]['content']) else search_result[search_result_index]['htmlText']
            search_result[search_result_index]['answer_for_question'] = search_content
            all_search_result.append(search_result[search_result_index]) 
        return all_search_result
    
    
    def messages_to_str(self,messages):
        messages_str = ''
        for item in messages:
            messages_str+=item['role']+':\n'
            messages_str+=item['content']+'\n'
        return messages_str
    @staticmethod
    def get_temperature_list(result_num):
        """
        temperature 尽可能均匀分布
        """
        if result_num <= 1:
            return [0.01]
        elif result_num == 2:
            return [0.01, 0.99]
        else:
            step = (0.99 - 0.01) / (result_num - 1)
            return [round(0.01 + i * step, 2) for i in range(result_num)]

    def __call__(self, content):
        return self.process_text(content)

    
    def log(self, message: str, title: str = "Info"):
        print(f"{title}: {message}")