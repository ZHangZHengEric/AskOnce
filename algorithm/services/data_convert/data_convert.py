from typing import List
import requests
import json
from time import sleep
import argparse
import os,sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.realpath(__file__)))))))
from AskOnce.algorithm.services.task_manager.task_manager import TaskManager,http_post_json
from AskOnce.algorithm.lib.data_convert.factory import BasicLoaderFactory
import time
from datetime import datetime
from urllib.parse import urlparse, unquote
import traceback
import re

class DataInput:
    def __init__(self,json_data,task_id) -> None:
        self.file_path = json_data['file_path']  # 需要检查是否是url
        self.is_remove_wrap = json_data['remove_wrap'] if 'remove_wrap' in json_data.keys() else False
        # self.file_type = None if "file_type" not in json_data else json_data['file_type']
        self.start_page=(1 if json_data['start_page'] ==-1 else json_data['start_page']) if 'start_page' in json_data.keys() else 1
        self.end_page= (None if json_data['end_page'] ==-1 else json_data['end_page'] ) if 'start_page' in json_data.keys() else None
        self.task_id = task_id

# 将输入字符串解析到输入结构体中
def unmarshal_task_input(GetTaskResp : dict):
    task_type = GetTaskResp['task_type']
    input_json = json.loads(GetTaskResp["input"])
    if task_type ==args.tasktype[0]:
        if type(input_json)==list:
            return  task_type,[DataInput(json_data=input_item, task_id=GetTaskResp["task_id"]) for input_item in input_json] 
        else:
            return task_type,[DataInput(json_data=input_json, task_id=GetTaskResp["task_id"]) ]
    
url_regex = re.compile(
    r'^(?:http|ftp)s?://'  # http:// or https:// or ftp:// or ftps://
    r'(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|'  # domain...
    r'localhost|'  # localhost...
    r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})'  # ...or IP
    r'(?::\d+)?'  # optional port
    r'(?:/?|[/?]\S+)$', re.IGNORECASE
)

def dowmload_url_to_path(url):
    parsed_url = unquote(url)
    parsed_url = urlparse(parsed_url)
    local_file_path = os.path.join(os.environ.get('CONVERT_CACHE'),os.path.basename(parsed_url.path))
    print('下载到本地地址',local_file_path)
    start_download_time = time.time()
    response = requests.get(url)
    if response.status_code == 200:
        try:
            with open(local_file_path, 'wb') as file:
                file.write(response.content)
                print(f'文件已下载至：{local_file_path}')
        except:
            print()
    else:
        print('下载失败，状态码：', response.status_code)
    print(f"下载文件时间: {time.time() - start_download_time}")
    response.close()
    return local_file_path    

def is_url(string):
    """判断一个字符串是否是链接"""
    return re.match(url_regex, string) is not None

# 执行任务
def process(task_input,task_type,model,args,tm):
    if task_type ==args.tasktype[0]:
        start_time = time.time()
        list_result = []
        for data_item in task_input:
            result_all = {} 
            # 根据后端传入的url, 保存文件
            file_path = data_item.file_path
            print('file_path',file_path)
            # file_type = data_item.file_type  # 肯定不会有这个参数。但是可以保留
            # 首先要判断file_path 是不是要下载的url 
            if is_url(file_path):
                print('file path 是url')
                file_path = dowmload_url_to_path(file_path)
            elif os.path.exists(file_path)==False:
                print('file path 不存在')
                raise ValueError('{file_path} not exists'.format(file_path=file_path))
            try: 
                # text 还是全部的有价值文本，text_detail 是更加详细的结构化解析，但是pdf 得与之前一样。其他的我们这次可以加上一些。
                text,meta_data = model.create(file_path) 
                result_all['text_detail'] = meta_data
                result_all['text'] = text
                print('解析后字数长度：',len(result_all['text'])) 
                list_result.append(result_all)
            except:
                print(traceback.format_exc())
                result_all['file_type'] = None
                result_all['text_detail'] = None
                result_all['text'] = ''
                list_result.append(result_all)
                print(f"{file_path} 文件解析异常")
        end_time = time.time()
        print('文件解析','用时'+str(end_time-start_time))
        return list_result
    
if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--jobdurl", type=str, default='', help="jobd use url not use ip & port")
    parser.add_argument("--tasktype", type=str, nargs='+', default=[], help="task_type")
    parser.add_argument("--worker_name", type=str, default='', help="worker_name")
    parser.add_argument("--correct_dict_path", type=str, help="correct_dict_path")
    args = parser.parse_args()
    tm = TaskManager(jobdurl=args.jobdurl)
    model = BasicLoaderFactory(args.correct_dict_path)    
    for one_task_type in args.tasktype:
        tm.add_task_type_info(one_task_type,10000,args.worker_name)

    while (True):
        try:
            errcode, get_task_resp = tm.get_task_block(args.tasktype)
            print('收到任务',datetime.now())
        except Exception as e:
            if 'value' not in str(e):
                print('gettask',e)
            continue
        if errcode != 200:
            sleep(0.05)
            continue
        elif "code" in get_task_resp.keys():
            sleep(0.05)
            continue
        else:
            try:
                one_task_type,task_input = unmarshal_task_input(get_task_resp)
            except Exception as e:
                print('process input_data',e)
                info = traceback.format_exc()
                print(info)
                tm.update_info(task_info={"task_id": get_task_resp["task_id"],"output": 'error',"status" : TaskManager.STATUS_INPUT_MISMATCH})
                continue
            try:
                taskOutput = process(task_input,one_task_type,model,args,tm)
            except Exception as e:
                print('process',e)
                info = traceback.format_exc()
                print(info)
                tm.update_info(task_info={"task_id": get_task_resp["task_id"],"output": 'error',"status" : TaskManager.STATUS_EXEC_FAILED})
                continue
            tm.update_info(task_info={"task_id": get_task_resp["task_id"],"output": json.dumps(taskOutput,ensure_ascii=False),"status" : TaskManager.STATUS_FINISH})
       
    
        
        
    
        
        