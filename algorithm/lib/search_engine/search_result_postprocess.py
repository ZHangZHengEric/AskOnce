



class SearchResultPostProcess:
    def __init__(self):
        pass
    
    def proceess(self,search_results):
        search_results = self.rrf_fusion_for_search_chunks(search_results=search_results)
        search_results  = self.merge_overlap_chunk(search_results)
        return search_results
    def rrf_fusion_for_search_chunks(self,search_results,rrf_k = 1):
        
        search_source = set([item['source'] for item in search_results])
        search_results_by_source = {item:[] for item in search_source}
        document_nums_by_source = {item:{} for item in search_source}
        for result in search_results:
            search_results_by_source[result['source']].append(result)
            if result['doc_id'] not in document_nums_by_source[result['source']]:
                document_nums_by_source[result['source']][result['doc_id']]=0
            document_nums_by_source[result['source']][result['doc_id']]+=1
        
        
        search_results_new= []
        for source_item in search_results_by_source.keys():
            search_results_by_source[source_item] = sorted(search_results_by_source[source_item],key=lambda x :x['score'],reverse=True)
            max_score = max(item['score'] for item in search_results_by_source[source_item]  )
            min_score = min(item['score'] for item in search_results_by_source[source_item]  )
            for index in range(len(search_results_by_source[source_item])):
                search_results_by_source[source_item][index]['ranking'] = index+1
                search_results_by_source[source_item][index]['normal_score']=  (search_results_by_source[source_item][index]["score"] - min_score) / (max_score - min_score)
            search_results_new.extend(search_results_by_source[source_item])
        search_results = search_results_new
        # print(search_results)
        
        # 提取出search_results 中所有的方法
        merged_results = {}
        for result in search_results:
            chunk = result["doc_content"]
            if chunk not in merged_results:
                merged_results[chunk] = {
                    "scores": {},
                    "normal_scores":{},
                    "doc_id": result['doc_id'],
                    "rankings":{},
                    "doc_content":chunk,
                    "sources":[],
                    "doc_segment_id" : result['doc_segment_id'],
                    "start" : result['start'],
                    "end" : result['end']
                }
            merged_results[chunk]['sources'].append(result['source'])
            merged_results[chunk]['scores'][result['source']]=result['score']
            merged_results[chunk]['normal_scores'][result['source']]=result['normal_score']
            merged_results[chunk]['rankings'][result['source']]=result['ranking']
        
        for chunk in merged_results.keys():
            rrf_score = 0
            for search_source_item in search_source:
                rank = merged_results[chunk]["rankings"].get(search_source_item, len(search_results_by_source[search_source_item])+1)
                score = merged_results[chunk]['normal_scores'].get(search_source_item,0)
                document_frequency_in_source = document_nums_by_source[search_source_item].get(merged_results[chunk]['doc_id'],0)
                document_frequency =  len(merged_results[chunk]["rankings"].keys())  /len(search_source)
                adjusted_score = score * (1 + 0.05 * document_frequency_in_source)
                rrf_score += document_frequency * (adjusted_score / (rrf_k + rank))
            merged_results[chunk]["score"] = rrf_score
        
        sorted_merged_results = dict(sorted(merged_results.items(), key=lambda x: x[1]["score"], reverse=True))        
        return list(sorted_merged_results.values())
        
    def merge_overlap_chunk(self,search_results):
        new_search_result = []
        doc_dict = {}
        for result in search_results:
            doc_name = result["doc_id"]
            if doc_name not in doc_dict:
                doc_dict[doc_name] = []
            doc_dict[doc_name].append(result)
        for chunks in doc_dict.values():
            chunks.sort(key=lambda x: x["start"])

        for doc_id , doc_chunks in doc_dict.items():
            i = 0
            while i < len(doc_chunks):
                current_chunk = doc_chunks[i]
                merged_chunk = current_chunk.copy()
                merged_chunk['doc_segment_ids'] = str(merged_chunk['doc_segment_id'])
                j = i + 1
                while j < len(doc_chunks):
                    next_chunk = doc_chunks[j]
                    if next_chunk["start"] <= merged_chunk["end"]:
                        # 存在重叠，更新合并chunk的结束索引和分数（取最大分数）
                        merged_chunk["score"] = max(merged_chunk["score"], next_chunk["score"])
                        # if 28325 == doc_id:
                        #     print(merged_chunk['doc_content'])
                        #     print(merged_chunk['start'],merged_chunk['end'],next_chunk["start"],next_chunk["end"])
                        #     print(len(next_chunk["doc_content"]))
                        if merged_chunk["end"] - next_chunk["start"]>=0:
                            merged_chunk['doc_content'] += next_chunk["doc_content"][merged_chunk["end"] - next_chunk["start"]:]
                        merged_chunk["end"] = max(merged_chunk["end"], next_chunk["end"])

                        # if 28325 == doc_id:
                        #     print(merged_chunk['doc_content'])
                            
                        merged_chunk['doc_segment_ids']+='+'+str(next_chunk['doc_segment_id'])
                        j += 1
                    else:
                        break
                new_search_result.append(merged_chunk)
                i = j
        # print(new_search_result)
        return new_search_result
        
        
    def remove_redundant(
            self,
            result_list,
            merge_length=5120, 
            merge_times=10,
            group_key="doc_id",
            merge_key="passage_content",
            sorted_key="passage_bm25_score"
        ):
        
        stop_merge = {'merge_length':merge_length,'merge_times':merge_times}
        # 先按照doc_name分组
        group_by_doc_name = {}
        for item in result_list:
            if item[group_key] not in group_by_doc_name.keys():
                group_by_doc_name[item[group_key]] = []
            group_by_doc_name[item[group_key]].append(item)
        
        # 对小组成员按照start升序
        result_list = []
        for key, group in group_by_doc_name.items():
            group_list = sorted(group, key=lambda x:x['start'])
            result_list.extend(group_list)
        
        key_check = result_list[0]
        if "start" not in key_check.keys() and  "end" not in key_check.keys():
            raise Exception("key start 或者 key end 不在元素中")
        
        new_result_list = []
        for i, result in enumerate(result_list):
            
            if i == 0 or \
            result[group_key] != new_result_list[-1][group_key] or \
            result['start'] > new_result_list[-1]['end'] or \
            new_result_list[-1]['text_length'] >= stop_merge['merge_length'] or \
            new_result_list[-1]['merge_times'] >= stop_merge['merge_times']:
                result['merge_times'] = 0
                result['text_length'] = result['end'] - result['start'] + 1
                new_result_list.append(result)
            else:
                if result['end'] > new_result_list[-1]['end']:
                    new_result_list[-1][merge_key] += result[merge_key][-(result["end"] - new_result_list[-1]["end"]):]
                    new_result_list[-1]['end'] = result['end']
                    new_result_list[-1]['text_length'] = new_result_list[-1]['end'] - new_result_list[-1]['start'] + 1
                    new_result_list[-1]["merge_times"] += 1
                    new_result_list[-1][sorted_key] = max(new_result_list[-1][sorted_key], result[sorted_key])
                else:
                    new_result_list[-1][sorted_key] = max(new_result_list[-1][sorted_key], result[sorted_key])
                    new_result_list[-1]['merge_times'] += 1 # 完全cover情况
        
        # 去除重复文字
        merge_results = []
        unique_list = []
        for item in new_result_list:
            if item[merge_key] not in unique_list:
                unique_list.append(item[merge_key])
                merge_results.append(item)
        merge_results = sorted(merge_results,key=lambda x: x[sorted_key], reverse=True)
        return merge_results
