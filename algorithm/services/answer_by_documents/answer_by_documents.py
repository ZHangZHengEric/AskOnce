from typing import List
import requests
import json
from time import sleep
import argparse
import os,sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.realpath(__file__)))))))
from AskOnce.algorithm.services.task_manager.task_manager import TaskManager,http_post_json
from AskOnce.algorithm.lib.llm_api.search_answer import SearchAnswer
import time
from datetime import datetime
import traceback

class QAnswerInput:
    def __init__(self,json_data,task_id) -> None:
        self.question = json_data['question']
        self.answer_style = json_data['answer_style'] if 'answer_style' in json_data.keys() else 'simplify'
        self.search_result = json_data['search_result'] if 'search_result' in json_data.keys() else []
        self.is_stream = json_data['is_stream'] if 'is_stream' in json_data.keys() else False
        self.task_id = task_id

# 将输入字符串解析到输入结构体中
def unmarshal_task_input(GetTaskResp : dict):
    # log_f = open(log_txt, 'a')
    task_type = GetTaskResp['task_type']
    # try:
    #     log_f.write(task_type+'\t'+str(time.strftime('%Y-%m-%d %H:%M:%S' ,time.localtime(time.time()))  )+'\t'+GetTaskResp["input"]+'\n')
    # except:
    #     pass
    # log_f.close()
    input_json = json.loads(GetTaskResp["input"])
    if task_type ==args.tasktype[0]:
        return  task_type,QAnswerInput(json_data=input_json, task_id=GetTaskResp["task_id"]) 
# 执行任务
def process(task_input,task_type,model,args,tm):
    if task_type ==args.tasktype[0]:
        start_time = time.time()
        result_all = {'answer':''} 
        print(task_input.answer_style)
        if task_input.answer_style =='simplify':
            result  = model.simplify_answer(task_input.question,task_input.search_result,task_input.is_stream)
        elif task_input.answer_style == 'detailed':
            result = model.detailed_answer(task_input.question,task_input.search_result,task_input.is_stream)
        elif task_input.answer_style == 'detailed_no_chapter':
            result = model.detailed_no_chapter_answer(task_input.question,task_input.search_result,task_input.is_stream)
        elif task_input.answer_style == 'askProfessional':
            result = model.professional_answer(task_input.question,task_input.search_result,task_input.is_stream)
        else:
            result = model.simplify_answer(task_input.question,task_input.search_result,task_input.is_stream)
        
        if task_input.is_stream:
            is_update = 0
            for content_part in result:
                result_all['answer'] += content_part
                try:
                    # print(result_all)
                    if is_update % 4==0: 
                        # up_start_time = time.time()
                        tm.update_info(task_info={"task_id": task_input.task_id,"output": json.dumps(result_all,ensure_ascii=False),"status" : TaskManager.STATUS_RUNNING},is_multi_resps=True)
                        # end_start_time = time.time()
                        # update_time += (end_start_time-up_start_time)
                        is_update=0
                    is_update+=1
                except:
                    print(traceback.format_exc())
                    print('流式返回出现错误')
        else:
            result_all['answer']= result
        end_time = time.time()
        print('回答','用时'+str(end_time-start_time))
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
    model = SearchAnswer(api_key=args.api_key,platform_api_url=args.platform_api_url,model_name= args.model_name)    
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
       
    
        
        
    
        
        