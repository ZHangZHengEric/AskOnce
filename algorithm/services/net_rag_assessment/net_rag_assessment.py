from typing import List
import requests
import json
from time import sleep
import argparse
import os,sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.realpath(__file__)))))))
from AskOnce.algorithm.services.task_manager.task_manager import TaskManager,http_post_json
from AskOnce.algorithm.lib.llm_api.question_process import QuestionProcess
import time
from datetime import datetime
import traceback
# history is a list[{'role':'user','message':""},{'role':'assistant','message':""},]
class DataInput:
    def __init__(self,json_data,task_id) -> None:
        self.question = json_data['question']
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
        result_all['result']= model.judge_use_rag(task_input.question)
        print(task_input.question,result_all['result'])
        end_time = time.time()
        print('判断是否需要互联网rag','用时'+str(end_time-start_time))
        return result_all
    
if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--jobdurl", type=str, default='', help="jobd use url not use ip & port")
    parser.add_argument("--tasktype", type=str, nargs='+', default=[], help="task_type")
    parser.add_argument("--worker_name", type=str, default='', help="worker_name")
    parser.add_argument("--api_key", type=str, help="api_key")
    parser.add_argument("--model_name", type=str, help="model_name")
    parser.add_argument("--platform_api_url", type=str, help="platform_api_url")
    # parser.add_argument("--log_file_path", type=str, default="/data1/zhangzheng/online_log/entity_normalization.log")
    args = parser.parse_args()
    # log_txt = args.log_file_path[:-3]+'input_log' 
    tm = TaskManager(jobdurl=args.jobdurl)
    model = QuestionProcess(api_key=args.api_key,platform_api_url=args.platform_api_url,model_name= args.model_name)    
    for one_task_type in args.tasktype:
        # 
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
       
    
        
        
    
        
        