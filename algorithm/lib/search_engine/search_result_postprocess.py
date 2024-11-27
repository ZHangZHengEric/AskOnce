def remove_redundant(
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
