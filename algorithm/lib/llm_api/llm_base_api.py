#!/usr/bin/env python
# -*-coding:utf-8 -*-
# @author: zhangjikang 
# 文本处理的基类, 其他几个文本生成类都继承自这个类
import time
from typing import Any

from openai import OpenAI


class LLMBaseAPI:
    def __init__(
            self, 
            platform_api_url: str , 
            api_key: str ,
            model_name: str 
        ):
        self.client = OpenAI(
            api_key=api_key,
            base_url=platform_api_url
        )
        self.model_name = model_name
    
    def ask_llm(
            self, 
            prompt: str | list, 
            temperature: float = 0.01, 
            presence_penalty: float = 1.2, 
            frequency_penalty: float = 1.2,
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
            max_tokens=max_tokens
        ) 
        return_content = completion.choices[0].message.content
        print('使用模型', self.model_name, '推理耗时', f'{time.time() - start_time:.4f} 秒')
        return return_content
    
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