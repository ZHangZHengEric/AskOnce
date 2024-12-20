from typing import List
import requests
import json
from time import sleep
import argparse
import os,sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.realpath(__file__)))))))
from AskOnce.algorithm.services.task_manager.task_manager import TaskManager,http_post_json
from AskOnce.algorithm.lib.search_engine.document_split import TextChunkPorcessor
import time
from datetime import datetime
import traceback

class DataInput:
    def __init__(self,json_data,task_id) -> None:
        self.text = json_data['text']
        self.text_cutting_version ="punc_cutting" if "text_cutting_version" not in json_data.keys() else json_data['text_cutting_version']
        self.window_size = 256 if 'window_size' not in json_data.keys() or json_data['window_size'] ==0 else int(json_data['window_size'])
        self.stride = 170 if 'stride' not in json_data.keys() or json_data['stride'] ==0 else int(json_data['stride'])
        self.fix_length_list = [128, 256, 512] if 'fix_length_list' not in json_data.keys() or len(json_data['fix_length_list']) == 0 else list(json_data['fix_length_list'])
        self.task_id = task_id

# 将输入字符串解析到输入结构体中
def unmarshal_task_input(GetTaskResp : dict):
    task_type = GetTaskResp['task_type']
    input_json = json.loads(GetTaskResp["input"])
    if task_type ==args.tasktype[0]:
        return  task_type,DataInput(json_data=input_json, task_id=GetTaskResp["task_id"]) 
# 执行任务
def process(task_input,task_type,model,args,tm):
    if task_type ==args.tasktype[0]:
        start_time = time.time()
        result_all = {} 
        
        if task_input.text_cutting_version == 'punc_cutting':
            sentences_list = model.merge_sentences_split(task_input.text,task_input.fix_length_list)
            result_all['sentences_list'] = sentences_list
        
        elif task_input.text_cutting_version == "move_window_cutting":
            sentences_list = model.move_window_split(task_input.text,window_size = task_input.window_size, stride= task_input.stride)
            result_all['sentences_list']  = sentences_list
        else:
            raise Exception("text_cutting_version not in ['punc_cutting','move_window_cutting']")   
        end_time = time.time()
        print('切分','用时'+str(end_time-start_time),'切割后句子数',len(result_all['sentences_list']))
        return result_all
    
if __name__ == '__main__':
    default_jobd_url = os.getenv('JOBD_URL', '')
    print('读取环境变量jobd_url:',default_jobd_url)
    parser = argparse.ArgumentParser()
    parser.add_argument("--jobdurl", type=str, default=default_jobd_url, help="jobd use url not use ip & port")
    parser.add_argument("--tasktype", type=str, nargs='+', default=[], help="task_type")
    parser.add_argument("--worker_name", type=str, default='', help="worker_name")
    args = parser.parse_args()
    tm = TaskManager(jobdurl=args.jobdurl)
    model = TextChunkPorcessor()
    for one_task_type in args.tasktype:
        tm.add_task_type_info(one_task_type,10000,args.worker_name)

    while (True):
        try:
            errcode, get_task_resp = tm.get_task_block(args.tasktype)
            if 'data' in get_task_resp.keys():
                get_task_resp = get_task_resp['data']
        except Exception as e:
            if 'value' not in str(e):
                print('gettask',e)
            continue
        if errcode != 200:
            sleep(0.05)
            print('空任务队列',datetime.now(),errcode)
            continue
        else:
            print('收到任务',datetime.now())
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
       
    
        
        
    
        
        