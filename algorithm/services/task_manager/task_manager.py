import requests
import json
from time import sleep
import os
import traceback
os.environ['CURL_CA_BUNDLE'] = ''

def http_post_json(url : str, body : dict):
    # try:
    resp = requests.post(url=url, json=body)
    # print(resp)
    # print('resp.text:',resp.text)
    try:
        body = json.loads(resp.text)
        resp.close()
        return resp.status_code, body
    except:
        traceback.print_exc()
        print('error',body,resp)
        resp.close()
        return 502,''
    

class TaskManager:
    # status of task
    STATUS_ENQUE_ERR   = "ENQUE_ERROR"
    STATUS_FINISH      = "FINISH"
    STATUS_CANCEL      = "CANCEL"
    STATUS_INPUT_MISMATCH = 'INPUT_MISMATCH'
    STATUS_EXEC_FAILED = "EXEC_FAILED"
    STATUS_WAITTING    = "WAITTING"
    STATUS_RUNNING     = "RUNNING"
    def __init__(self,jobdip='127.0.0.1',jobdport=20033,jobdurl='',worker_name='') -> None:
        self.worker_name = worker_name
        if len(jobdurl)==0:
            default_jobd_url = os.getenv('JOBD_ADDR', '')
            print('读取环境变量jobd_url:',default_jobd_url)
            self.jobdurl = default_jobd_url
        else:
            self.jobdurl = jobdurl
        self.GetTaskUrl = "{}/jobd/worker/GetTask"
        self.BlockBatchGetTaskUrl = "{}/jobd/worker/BlockBatchGetTask"
        self.UpdateInfoUrl = "{}/jobd/worker/UpdateInfo"
        self.ClearTaskUrl = '{}/jobd/api/ClearTask'
        self.AddTaskTypeInfoUrl = '{}/jobd/api/AddTaskTypeInfo'
        self.UpdateTaskTypeInfoUrl = '{}/jobd/api/UpdateTaskTypeInfo'
        self.GetTaskNumUrl = '{}/jobd/api/GetTaskNum'
        self.DoTaskUrl = "{}/jobd/committer/DoTask"
    # 清除某个任务队列
    def clear_task(self,task_type):
        url = self.ClearTaskUrl.format(self.jobdurl)
        body = {
            "task_type": task_type
        }
        errcode, resp = http_post_json(url, body)
        print(errcode, resp)
        return errcode,resp
    
    # 添加某个任务类型的队列
    def add_task_type_info(self,task_type,task_num_limit,worker_name=''):
        url = self.AddTaskTypeInfoUrl.format(self.jobdurl)
        body = {
            "task_type": task_type,
            "task_num_limit": task_num_limit,
            "instance":self.worker_name
        }
        errcode, resp = http_post_json(url, body)
        print(errcode, resp)
        return errcode,resp

    # 获取一个新的任务
    def get_task(self,task_type):
        errcode, resp = http_post_json(url=self.GetTaskUrl.format(self.jobdurl), body={
            "task_type": task_type,
            "instance":self.worker_name
        })
        # print("get task:", errcode, resp)
        return errcode, resp

    # 获取一个新的任务
    def get_task_block(self,task_types):
        errcode, resp = http_post_json(url=self.BlockBatchGetTaskUrl.format(self.jobdurl), body={
            "task_types": task_types,
            "instance":self.worker_name
        })
        # print("get task:", errcode, resp)
        return errcode, resp

    # 上传任务结果，或者更新任务进度
    def update_info(self,task_info : dict,is_multi_resps=False) -> int:
        if is_multi_resps==False or task_info['status']!=self.STATUS_RUNNING:
            print('start update info')
        errcode, resp = http_post_json(url=self.UpdateInfoUrl.format(self.jobdurl), body={
            "task_id": task_info['task_id'],
            "output": task_info['output'],
            "instance":self.worker_name,
            "status": task_info['status'],
            "extend_info" : { 
                "is_multi_resps" : is_multi_resps,  # 这里需要设置成True
            }
        })
        if is_multi_resps==False or task_info['status']!=self.STATUS_RUNNING:
            print("update info:", errcode, resp)
        return errcode
    def update_task_type_info(self,task_type,task_num_limit):
        url = self.UpdateTaskTypeInfoUrl.format(self.jobdurl)
        body = {
            "task_type": task_type,
            "task_num_limit": task_num_limit
        }
        errcode, resp = http_post_json(url, body)
        print(errcode, resp)
    def get_task_num(self,task_type):
        url = self.GetTaskNumUrl.format(self.jobdurl)
        body = {
            "task_type": task_type,
        }
        errcode, resp = http_post_json(url, body)
        print(errcode, resp)

    def do_task(self,task_type,input_data,timeout=30*1000):
        
        errcode, resp = http_post_json(self.DoTaskUrl.format(self.jobdurl),{
            "task_type": task_type,
            "input": input_data if type(input_data) == str else json.dumps(input_data,ensure_ascii=False),
            "timeout_ms": timeout,
        })
        print(errcode, resp)
        return errcode,resp