

# 一个通用的es8 curd 服务，目前在支持线上askonce

import os
import time
import logging
import argparse
import json
import pandas as pd
from glob import glob
from tqdm import tqdm
from datetime import datetime
from typing import List, Dict, Any
from pathlib import Path
from elasticsearch.helpers import bulk, parallel_bulk
from elasticsearch import Elasticsearch
from elasticsearch.exceptions import NotFoundError

# from AtomES.lib.atom_utils.common_utils import genenrate_id_based_string

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
    datefmt="%m/%d/%Y %H:%M:%S",
)


# os.environ['ES_SETTING'] = '/mnt/nvme0n1/chendong/aiwork/code/AtomES/lib/db_model/db_config/es_setting.json'

logger = logging.getLogger(__name__)

class Es8Config:
    ''' 
    es8 的相关配置，可以参考写的一个例子 ：
        /mnt/nvme0n1/chendong/aiwork/code/AtomES/lib/db_model/db_config/es_config_example.json
    这个mapping_obj__or_path 后续由后端入参传入
    '''
    @staticmethod
    def load_mapping_config(mapping_obj_or_path):
        if isinstance(mapping_obj_or_path, Dict):
            mapping_config = mapping_obj_or_path
            return mapping_config
        else:
            mapping_config = json.load(open(mapping_obj_or_path))
            print("mapping config:", mapping_config)
            return mapping_config

class Es8Interface:
    
    def __init__(self,address,username,password, es_setting_path) -> None:
        # self.mapper = mapper
        # self.index =self.mapper['mappings_config']['index_name']
        # /home/nlp_platform/AtomES/lib/db_model/db_config/es_config_example.json
        # self.setting = JSONSerializable.load_from_file(os.environ['ES_SETTING'])
        # script_directory = 
        # file_path = Path(__file__).parent / 'db_config' / 'es_setting.json'
        self.setting = json.load(open(es_setting_path))
        print("获取es8配置成功")
        print(self.setting)
        try:
            self.es = Elasticsearch(
                [address], bearer_auth=(username, password), request_timeout=30,
            )
            logger.info("es engine 初始化成功")
        except:
            raise Exception("es engine 初始化出错")
    
    def conn(self,mapper):
        index_name = mapper['mappings_config']['index_name']
        try:
            if self.es.indices.exists(index=index_name):
                logger.info(f"The index {index_name} exists, es conn 成功")
            else:
                logger.info(f"The index {index_name} does not exist,start to creating...")
                self.create(mapper)
        except:
            raise Exception("连接数据库失败")
    
    def create(self,mapper):
        '''和业务无关, 只和 config 相干，一次性创建好业务所有需要的index'''
        settings = self.setting # es node 级别
        mapper_config = mapper # List, 可能一次性需要创建多个index
        index_name = mapper['mappings_config']['index_name']
        index_mappings = mapper_config['mappings_config']['mappings']
        try:
            create_info = self.es.indices.create(index=index_name, mappings=index_mappings, settings=settings['settings'])
            logger.info(create_info)
        except:
            raise Exception("创建数据库失败")
            
    def insert(self, texts,mapper):
        index_name = mapper['mappings_config']['index_name']
        insert_results = []
        print(texts[0])
        insert_report = parallel_bulk(self.es, texts, index=index_name)
        for rep in insert_report:
            if rep[0] ==False:
                insert_results.append(rep[1]['index']['doc_id'])
        print(insert_results)
        return insert_results
    
    def query_analysis(self, query,mapper):
        '''
            to do: ES 分词替换
        '''
        index_name = mapper['mappings_config']['index_name']
        query_token = []
        response = self.es.indices.analyze(index=index_name, analyzer="my_ana", body={"text": query})
        tokens = response["tokens"]
        for token in tokens:
            query_token.append({"token": token['token'], "position":token['position']})
        return (query, query_token)
    
    
    def delete(self, doc_id_list:List, mapper,delete_all = False,):
        '''
            delete_all = True 的情况下, 删除index, 否则 删除 指定的doc_id_list
        '''
        if delete_all:
            try:
                # 尝试删除索引
                index_name = mapper['mappings_config']['index_name']
                response = self.es.indices.delete(index=index_name)
                logger.info(f"索引 '{index_name}' 删除成功: {response}")
                return []
            except NotFoundError:
                logger.error(f"索引 '{index_name}' 不存在，无法删除。")
            except Exception as e:
                logger.error(f"删除索引时出错: {e}")
        else:    
            failed_doc_ids = []
            if not isinstance(doc_id_list, List):
                raise TypeError("doc_id_list must be list type")
            for doc_id in doc_id_list:
                body = {
                    "query": {
                        "match": {
                            "doc_id": doc_id
                        }
                    }
                }
                try:
                    response = self.es.delete_by_query(index=mapper['mappings_config']['index_name'], body=body)
                    logger.info(f'Deletion {doc_id} operation successful.')
                except Exception as e:
                    logger.info(f'Deletion operation failed: {str(e)}')
                    for failure in response.get('failures', []):
                        doc_id = failure.get('_source', {}).get('doc_id')
                    if doc_id:
                        failed_doc_ids.append(doc_id)
            # 返回删除失败的文档ID列表
            return failed_doc_ids

class BasicEs8(Es8Interface):
    def __init__(self, address,username,password,es_setting_path) -> None:
        super().__init__(address,username,password,es_setting_path)
    
    def bm25_search(self, search_body,mapper ,return_fields):
        index_name = mapper['mappings_config']['index_name']
        results = []
        bm25_search_result = self.es.search(index=index_name,body=search_body)
        # query = search_body['query']['match']['doc_content']
        # query_seg_result = self.query_analysis(query)
        # print(query_seg_result)
        hits = bm25_search_result["hits"]["hits"]
        if len(hits)<=0: return results
        logger.info(f"return_fields: {return_fields}")
        if len(return_fields) == 0 and len(hits) > 0: return_fields = hits[0]['_source'].keys()
        if len(hits) > 0:
            for hit in hits:
                results.append(
                    {
                        "source" : {k : v  for k, v in hit["_source"].items() if k in return_fields},
                        "score": hit["_score"],                         
                    }
                )
        return results
    
    def emb_search(self, search_body, mapper,return_fields):
        # 执行向量检索查询
        results = []
        index_name = mapper['mappings_config']['index_name']
        emb_search_result = self.es.search(index=index_name,body=search_body)
        hits = emb_search_result['hits']['hits']
        print(hits)
        if len(hits) <= 0: return results
        if len(return_fields) == 0: return_fields = hits[0]['_source'].keys()
        # logger.info(f"return_fields: {return_fields}")
        for hit in hits:
            results.append(
                {
                    "source" : {k : v  for k, v in hit["_source"].items() if k in return_fields},
                    "score": hit["_score"],                         
                }
            )
        return results