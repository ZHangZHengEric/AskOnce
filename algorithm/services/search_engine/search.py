from typing import List
import requests
import json
from time import sleep
import argparse
import os,sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.realpath(__file__)))))))
from AskOnce.algorithm.services.task_manager.task_manager import TaskManager,http_post_json
from AskOnce.algorithm.lib.search_engine.es8_db_interface import BasicEs8
import time
from datetime import datetime
import traceback

class DataInput:
    def __init__(self,json_data,task_id) -> None:
        self.search_body = json_data['search_body'] # list类型
        self.search_type = json_data['search_type'] # 可选bm25, vec, all
        self.return_fields = []  if "return_fields" not in json_data else json_data["return_fields"]
        self.mapper_config = json_data['mapper_config']
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
        model.conn(task_input.mapper_config)        
        try:
            start_time = time.time()
            search_body = task_input.search_body
            search_type = task_input.search_type
            return_fields = task_input.return_fields
            if return_fields is None:
                return_fields = []

            if search_type == "vec":
                search_results = model.emb_search(search_body[0],task_input.mapper_config, return_fields)
            elif search_type == "bm25":
                search_results = model.bm25_search(search_body[0],task_input.mapper_config, return_fields)
            elif search_type == "all":
                if 'knn' in search_body[0].keys():
                    vec_search_results = model.emb_search(search_body[0],task_input.mapper_config,return_fields)
                    bm25_search_results = model.bm25_search(search_body[1],task_input.mapper_config,return_fields)
                    search_results = vec_search_results + bm25_search_results
                else:
                    vec_search_results = model.emb_search(search_body[1],task_input.mapper_config,return_fields)
                    bm25_search_results = model.bm25_search(search_body[0],task_input.mapper_config,return_fields)
                    search_results = vec_search_results + bm25_search_results
            print(f"es8检索耗时：{(time.time() - start_time)} ms")
            return search_results
        except:
            raise Exception("es8 检索报错")
    
if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--jobdurl", type=str, default='', help="jobd use url not use ip & port")
    parser.add_argument("--tasktype", type=str, nargs='+', default=[], help="task_type")
    parser.add_argument("--worker_name", type=str, default='', help="worker_name")
    parser.add_argument("--es_address", type=str, help="es_address")
    parser.add_argument("--es_username", type=str, help="es_username")
    parser.add_argument("--es_password", type=str, help="es_password")
    parser.add_argument("--es_setting_path", type=str, help="es_setting_path")
    args = parser.parse_args()
    tm = TaskManager(jobdurl=args.jobdurl)
    model = BasicEs8(args.es_address,args.es_username,args.es_password,args.es_setting_path)

    for one_task_type in args.tasktype:
        tm.add_task_type_info(one_task_type,10000,args.worker_name)

    while (True):
        try:
            errcode, get_task_resp = tm.get_task_block(args.tasktype)
            get_task_resp = get_task_resp['data']
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
       
    
        
        
    
        
        