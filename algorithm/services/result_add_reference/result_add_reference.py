from typing import List
import requests
import json
from time import sleep
import argparse
import os,sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.realpath(__file__)))))))
from AskOnce.algorithm.services.task_manager.task_manager import TaskManager,http_post_json
from AskOnce.algorithm.lib.result_reference.result_add_reference import AddReference
import time
from datetime import datetime
import traceback

class DataInput:
    def __init__(self,json_data,task_id) -> None:
        self.result = json_data['result']
        self.reference_list  = json_data['reference_list']
        self.threshold = json_data['threshold'] if 'threshold' in json_data.keys() else -1
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
        reference_map,reference_list_accept_index_range,reference_result = model.get_answer_with_reference_map(task_input.result,task_input.reference_list,task_input.threshold)
        result_all['reference_map'] = reference_map
        result_all['reference_result'] = reference_result
        result_all['reference_list_select_index'] = reference_list_accept_index_range
        print(result_all)
        end_time = time.time()
        print('匹配引用','用时'+str(end_time-start_time))
        return result_all
    
if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--jobdurl", type=str, default='', help="jobd use url not use ip & port")
    parser.add_argument("--tasktype", type=str, nargs='+', default=[], help="task_type")
    parser.add_argument("--worker_name", type=str, default='', help="worker_name")
    parser.add_argument("--stop_word_file", type=str, help="stop_word_file")
    args = parser.parse_args()
    tm = TaskManager(jobdurl=args.jobdurl)
    model = AddReference(stop_word_file = args.stop_word_file)    
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
       
    
        
        
    
        
        