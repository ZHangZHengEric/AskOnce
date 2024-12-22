import re
import difflib
import jieba

import logging
logging.basicConfig(level = logging.DEBUG,
                    format='%(asctime)s - %(levelname)s - %(message)s')
logger  = logging.getLogger(__name__)
class AddReference:
    def __init__(self,stop_word_file,threshold=0.5) -> None:
        self.stop_words = self.load_stop_words(stop_word_file)
        self.threshold = threshold
        self.init_threshold = threshold
        
    def load_stop_words(self,stop_words_file):
        f = open(stop_words_file)
        words = f.readlines()
        return words
    
    def split_answer(self,answer):
        all_parts =[]
        answer_parts  = re.split(r'([\n])',answer)
        for one_answer_part in answer_parts:
            if one_answer_part=='\n':
                all_parts.append('\n')
            elif len(one_answer_part)>0:
                para = one_answer_part
                para = re.sub('([•。;；！？\?])([^”’])', r"\1\n\2", para)  # 单字符断句符
                para = re.sub('(\.{6})([^”’])', r"\1\n\2", para)  # 英文省略号
                para = re.sub('(\…{2})([^”’])', r"\1\n\2", para)  # 中文省略号
                para = re.sub('([，,。！？\?][”’])([^，,。！？\?])', r'\1\n\2', para)
                # 如果双引号前有终止符，那么双引号才是句子的终点，把分句符\n放到双引号后，注意前面的几句都小心保留了双引号
                # para = para.rstrip()  # 段尾如果有多余的\n就去掉它
                # 很多规则中会考虑分号;，但是这里我把它忽略不计，破折号、英文双引号等同样忽略，需要的再做些简单调整即可。
                all_parts.extend(para.split("\n"))
        return all_parts


    def remove_stop_words_length(self,sentence):
        seg_list=list(jieba.cut(sentence, cut_all=False)) #精确模式
        # seg_list_index = 
        seg_list_not_stop_word= []
        seg_list_not_stop_word_length = 0
        for seg_item in seg_list:
            if seg_item not in self.stop_words:
                seg_list_not_stop_word.append(seg_item)
                seg_list_not_stop_word_length +=len(seg_item)
        return seg_list,seg_list_not_stop_word,seg_list_not_stop_word_length
    # 计算出回答结果片段，在原文中能够找到的比例。需要对回答结果偏进行分词，至少连续3个字符以上的字一直才算找到了
    # 判断sentence 找到reference_item一样的表述  在sentence 占比是多少
    def max_sub_sentence(self,sentence,reference_item):
        if len(reference_item)<10:
            return 0
        if len(sentence)<10:
            return 0
        seg_list,seg_list_not_stop_word,seg_list_not_stop_word_length = self.remove_stop_words_length(sentence)
        
        s = difflib.SequenceMatcher(None, sentence, reference_item)
        matches = s.get_matching_blocks()
        
        match_parts = []
        match_parts_length = 0
        for match in matches:
            output=sentence[match.a:match.a + match.size]
            #只显示大于1个字符的,要判断是中英文 英文要大于3个字符。 
            if len(output)>1:
                match_parts.append(output)
                # 计算长度的时候要去掉停用词
                output_seg_list,output_seg_list_not_stop_word,output_seg_list_not_stop_word_length = self.remove_stop_words_length(output)
                
                match_parts_length+= output_seg_list_not_stop_word_length
        
        return match_parts_length/(seg_list_not_stop_word_length+0.1)

    # 判断sentence 找到reference_item一样的表述  在 reference_item 占比是多少
    def max_sub_sentence_for_reference(self,sentence,reference_item):
        if len(reference_item)<5:
            return 0
        if len(sentence)<5:
            return 0
        seg_list,seg_list_not_stop_word,seg_list_not_stop_word_length = self.remove_stop_words_length(reference_item)
        
        s = difflib.SequenceMatcher(None, sentence, reference_item)
        matches = s.get_matching_blocks()
        
        match_parts = []
        match_parts_length = 0
        for match in matches:
            output=sentence[match.a:match.a + match.size]
            #只显示大于1个字符的,要判断是中英文 英文要大于3个字符。 
            if len(output)>1:
                match_parts.append(output)
                # 计算长度的时候要去掉停用词
                output_seg_list,output_seg_list_not_stop_word,output_seg_list_not_stop_word_length = self.remove_stop_words_length(output)
                match_parts_length+= output_seg_list_not_stop_word_length
        
        return match_parts_length/(seg_list_not_stop_word_length+0.1)
    
    # 找到回答 分割后的 的一个句子 与哪些引用一致
    def find_sentence_reference(self,sentence,reference_list):
        reference_list_item_index = []
        for reference_list_index , reference_list_item in enumerate(reference_list):
            ratio = self.max_sub_sentence(sentence,reference_list_item)
            logger.debug(f'{ratio},{sentence},||||,{reference_list_item}') 
            if ratio > self.threshold:
                reference_list_item_index.append(reference_list_index)
        return reference_list_item_index

    def find_use_reference(self,answer_part,reference_item):
        reference_item_parts = self.split_answer(reference_item)
        reference_item_parts_use_indexs = []
        reference_item_parts_index = []
        for reference_item_part in reference_item_parts:
            if len(reference_item_parts_index)==0:
                reference_item_parts_index.append([0,len(reference_item_part)])
            else:
                reference_item_parts_index.append([reference_item_parts_index[-1][1],reference_item_parts_index[-1][1]+len(reference_item_part)])
            ratio = self.max_sub_sentence_for_reference(reference_item_part,answer_part)
            if ratio > self.threshold:
                reference_item_parts_use_indexs.append(reference_item_parts_index[-1])
        if len(reference_item_parts_use_indexs)>0:
            return reference_item_parts_use_indexs
        else:
            return []
        
    def find_use_reference_new(self,answer_part,reference_item):
        reference_item_parts = self.split_answer(reference_item)
        reference_item_parts_use_indexs = []
        reference_item_parts_index = []
        for reference_item_part in reference_item_parts:
            if len(reference_item_parts_index)==0:
                reference_item_parts_index.append([0,len(reference_item_part)])
            else:
                reference_item_parts_index.append([reference_item_parts_index[-1][1],reference_item_parts_index[-1][1]+len(reference_item_part)])
            # 判断 reference_item_part 找到和answer_part一样的表述  在 answer_part 占比是多少
            ratio = self.max_sub_sentence_for_reference(reference_item_part,answer_part)
            if ratio > self.threshold:
                reference_item_parts_use_indexs.append({'index_range':reference_item_parts_index[-1],'ratio':ratio})
        if len(reference_item_parts_use_indexs)>0:
            return self.get_max_score_reference_part_index(reference_item_parts_use_indexs)
        else:
            return []
        
    def get_max_score_reference_part_index(self,reference_item_parts_use_indexs):
        if not reference_item_parts_use_indexs:
            return []
        reference_item_parts_use_indexs.sort(key=lambda x: x['index_range'][0])

        merged_indexs = []
        i = 0
        while i < len(reference_item_parts_use_indexs):
            current_index_range = reference_item_parts_use_indexs[i]['index_range']
            current_ratio = reference_item_parts_use_indexs[i]['ratio']
            end = current_index_range[1]
            j = i + 1
            while j < len(reference_item_parts_use_indexs) and reference_item_parts_use_indexs[j]['index_range'][0]-3 <= end:
                # 合并相邻区间，更新结束位置和取ratio最大值
                end = max(end, reference_item_parts_use_indexs[j]['index_range'][1])
                current_ratio = max(current_ratio, reference_item_parts_use_indexs[j]['ratio'])
                j += 1
            merged_indexs.append({'index_range': [current_index_range[0], end], 'ratio': current_ratio})
            i = j

        # 对合并后的结果按照ratio进行降序排序，获取ratio最大的index_range
        merged_indexs.sort(key=lambda x: x['ratio'], reverse=True)
        return merged_indexs[0]['index_range']

    
    def merge_and_drop_reference_index_range(self,index_ranges,merge_max_gap=15): 
        if len(index_ranges)==0:
            return []
        index_ranges = sorted(index_ranges,key = lambda x:x[0])
        new_index_ranges = []
        for index_ranges_item_index in range(len(index_ranges)):
            if len(new_index_ranges) ==0 :
                new_index_ranges.append(index_ranges[index_ranges_item_index])
            else:
                if index_ranges[index_ranges_item_index][0] - new_index_ranges[-1][1] < merge_max_gap:
                    new_index_ranges[-1][1] = index_ranges[index_ranges_item_index][1]
                else:
                    new_index_ranges.append(index_ranges[index_ranges_item_index])
        
        accept_index_range = [0,0]
        for new_index_ranges_item in new_index_ranges:
            if (new_index_ranges_item[1]- new_index_ranges_item[0]) > (accept_index_range[1]-accept_index_range[0]):
                accept_index_range = new_index_ranges_item
        return accept_index_range
            
    # 在给到答案以及 ，答案所引用的片段，找到与答案相对应的片段的精确部分。对reference_map 进行丰富。
    def find_part_index_from_reference(self,answer_parts,reference_list,reference_map):
        reference_list_accept_index_range = {}
        
        for answer_parts_index , item in enumerate(reference_map):
            item_reference_list = item['reference_list']
            # 对引用列表进行编列
            for item_reference_list_index in item_reference_list:
                # 找到 回答的chunk 对该引用文章中，去找对应chunk 的index
                use_range = self.find_use_reference( answer_parts[answer_parts_index],reference_list[item_reference_list_index])
                if item_reference_list_index not in reference_list_accept_index_range.keys():
                    reference_list_accept_index_range[item_reference_list_index] = []
                reference_list_accept_index_range[item_reference_list_index].extend(use_range)
        
        # 对reference_list_accept_index_range 进行合并调整
        for item_key in reference_list_accept_index_range.keys():
            reference_list_accept_index_range[item_key] = self.merge_and_drop_reference_index_range(reference_list_accept_index_range[item_key])
        return reference_list_accept_index_range
    
    def get_similar_chunk_range_from_refers(self,answer_parts,reference_list,reference_map_new):
        for one_result_chunk_index , one_result_chunk in enumerate(reference_map_new):
            answer_chunk = answer_parts[one_result_chunk_index]
            item_reference_list = one_result_chunk['refers']
            for item_reference_item_index,item_reference_item in enumerate(item_reference_list):
                use_range = self.find_use_reference_new( answer_chunk,reference_list[item_reference_item['index']])
                if len(use_range)>0:
                    reference_map_new[one_result_chunk_index]['refers'][item_reference_item_index]['referStart'] = use_range[0]
                    reference_map_new[one_result_chunk_index]['refers'][item_reference_item_index]['referEnd'] = use_range[1]
        return reference_map_new
    
    def merge_reference_map(self, reference_map):
        if not reference_map:
            return []
        # 先按照start位置对参考映射进行排序，便于后续查找相邻可合并的文本段
        reference_map.sort(key=lambda x: x['start'])
        merged_reference_map = []
        i = 0
        while i < len(reference_map):
            current_item = reference_map[i]
            current_start = current_item['start']
            current_end = current_item['end']
            current_refers = current_item['refers']
            j = i + 1
            while j < len(reference_map) and (reference_map[j]['start'] -3) <= current_end:
                # 检查引用文章是否一致，不一致则不能合并
                next_refers = reference_map[j]['refers']
                if self._check_refers_compatible(current_refers, next_refers):
                    current_end = max(current_end, reference_map[j]['end'])
                    current_refers = self._merge_refers(current_refers, next_refers)
                else:
                    break
                j += 1
            merged_reference_map.append({
                'start': current_start,
                'end': current_end,
                'refers': current_refers
            })
            i = j
        return merged_reference_map

    def _check_refers_compatible(self, refers1, refers2):
        """检查两组引用是否兼容（即引用文章是否一致）"""
        referred_articles1 = set([ref['index'] for ref in refers1])
        referred_articles2 = set([ref['index'] for ref in refers2])
        return referred_articles1 == referred_articles2
    
    def _merge_refers(self, refers1, refers2):
        """合并两组引用，确保每个文章只有一个对应片段"""
        merged_refers = {}
        all_refs = refers1 + refers2
        for ref in all_refs:
            index = ref['index']
            if index not in merged_refers:
                merged_refers[index] = ref.copy()
            else:
                merged_refers[index]['referStart'] = min(merged_refers[index]['referStart'], ref['referStart'])
                merged_refers[index]['referEnd'] = max(merged_refers[index]['referEnd'], ref['referEnd'])
        return list(merged_refers.values())
    
    # 在回答里面添加引用标记
    def get_answer_with_reference_map(self,answer,reference_list,threshold=-1):
        logger.info(f'输入的阈值:{threshold}')
        logger.info(f'reference list 长度:{len(reference_list)}')
        logger.info(f'answer的长度:{len(answer)}')
        if threshold > 0:
            self.threshold = threshold        
        answer_parts = self.split_answer(answer)
        logger.info(f'answer_parts 的长度:{sum([len(item) for item in answer_parts])}')
        if len(answer) != sum([len(item) for item in answer_parts]):
            logger.error(f'** 重大错误 **:{answer}')
        logger.info(answer_parts)
        # answer_parts_ref_indexs 是对应列表分割后的chunk 的 引用列表的index
        answer_parts_ref_indexs = []
        # answer_parts_index 放的是分割后的chunk 的index
        answer_parts_index = []
        for answer_part in answer_parts:
            
            if len(answer_parts_index)==0:
                answer_parts_index.append([0,len(answer_part)])
            else:
                answer_parts_index.append([answer_parts_index[-1][1],answer_parts_index[-1][1]+len(answer_part)])
            reference_indexs = self.find_sentence_reference(answer_part,reference_list=reference_list)
            answer_parts_ref_indexs.append(reference_indexs)
        logger.info(answer_parts_index)
        
        # 将answer_parts_index 和answer_parts_ref_indexs 放到reference_map 中
        reference_map = []
        reference_map_new = []
        for index in range(len(answer_parts)):
            reference_map.append({'index_range':answer_parts_index[index],'reference_list':answer_parts_ref_indexs[index]})
            reference_map_new.append({'start':answer_parts_index[index][0],'end':answer_parts_index[index][1],'refers':[{'index':item} for item in answer_parts_ref_indexs[index]]})
        
        # 找到chunk 在引用列表中的相似的chunk
        reference_list_accept_index_range = self.find_part_index_from_reference(answer_parts,reference_list,reference_map)
        reference_map_new = self.get_similar_chunk_range_from_refers(answer_parts,reference_list,reference_map_new)
        
        # 合并句子引用
        index = 0
        while True:
            if index >= len(reference_map)-1:
                break
            if len(reference_map[index]['reference_list'])==0:
                index+=1
                continue
            else:
                if len(reference_map[index+1]['reference_list'])==0:
                    index+=1
                    continue
                else:
                    if len(list(set(reference_map[index]['reference_list']) & set(reference_map[index+1]['reference_list'])))>0:
                        reference_map[index]['index_range'] = [reference_map[index]['index_range'][0],reference_map[index+1]['index_range'][1]]
                        reference_map[index]['reference_list'] = list(set(reference_map[index]['reference_list']).union(set(reference_map[index+1]['reference_list'])))
                        reference_map.pop(index+1)
                        index-=1
                    else:
                        index+=1
                        continue
            index+=1
        logger.info(reference_map)
        self.threshold = self.init_threshold
        
        # 对reference_map_new 进行合并
        reference_map_new= self.merge_reference_map(reference_map_new)
        logger.info(reference_map_new)
        
        return reference_map,reference_list_accept_index_range,reference_map_new