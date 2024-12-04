import re
import difflib
import jieba


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
                para = re.sub('([•。,，;；！？\?])([^”’])', r"\1\n\2", para)  # 单字符断句符
                para = re.sub('(\.{6})([^”’])', r"\1\n\2", para)  # 英文省略号
                para = re.sub('(\…{2})([^”’])', r"\1\n\2", para)  # 中文省略号
                para = re.sub('([。！？\?][”’])([^，。！？\?])', r'\1\n\2', para)
                # 如果双引号前有终止符，那么双引号才是句子的终点，把分句符\n放到双引号后，注意前面的几句都小心保留了双引号
                # para = para.rstrip()  # 段尾如果有多余的\n就去掉它
                # 很多规则中会考虑分号;，但是这里我把它忽略不计，破折号、英文双引号等同样忽略，需要的再做些简单调整即可。
                all_parts.extend(para.split("\n"))
        # return tokenizer.sent_tokenize(answer,language='chinese')
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
                # print(output)
                match_parts.append(output)
                # 计算长度的时候要去掉停用词
                output_seg_list,output_seg_list_not_stop_word,output_seg_list_not_stop_word_length = self.remove_stop_words_length(output)
                
                match_parts_length+= output_seg_list_not_stop_word_length
        
        return match_parts_length/(seg_list_not_stop_word_length+0.1)

    # 判断sentence 找到reference_item一样的表述  在 reference_item 占比是多少
    def max_sub_sentence_for_reference(self,sentence,reference_item):
        if len(reference_item)<5:
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
                # print(output)
                match_parts.append(output)
                # 计算长度的时候要去掉停用词
                output_seg_list,output_seg_list_not_stop_word,output_seg_list_not_stop_word_length = self.remove_stop_words_length(output)
                match_parts_length+= output_seg_list_not_stop_word_length
        
        return match_parts_length/(seg_list_not_stop_word_length+0.1)
    
    # 找到回答的一个句子与哪些引用一致
    def find_sentence_reference(self,sentence,reference_list):
        reference_list_item_index = []
        for reference_list_index , reference_list_item in enumerate(reference_list):
            ratio = self.max_sub_sentence(sentence,reference_list_item)
            print(ratio,sentence,reference_list_item) 
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
            # print(ratio,answer_part,reference_item_part)
            if ratio > self.threshold:
                reference_item_parts_use_indexs.append(reference_item_parts_index[-1])
        if len(reference_item_parts_use_indexs)>0:
            return reference_item_parts_use_indexs
        # [reference_item_parts_use_indexs[0][0],reference_item_parts_use_indexs[-1][1]]
        else:
            return []
    
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
            for item_reference_list_index in item_reference_list:
                # print(answer_parts[answer_parts_index],reference_list[item_reference_list_index])
                use_range = self.find_use_reference( answer_parts[answer_parts_index],reference_list[item_reference_list_index])
                if item_reference_list_index not in reference_list_accept_index_range.keys():
                    reference_list_accept_index_range[item_reference_list_index] = []
                reference_list_accept_index_range[item_reference_list_index].extend(use_range)
        
        # 对reference_list_accept_index_range 进行合并调整
        for item_key in reference_list_accept_index_range.keys():
            
            reference_list_accept_index_range[item_key] = self.merge_and_drop_reference_index_range(reference_list_accept_index_range[item_key])
            
            # range_min_index = len(reference_list[item_key])
            # range_max_index = 0
            # for item in reference_list_accept_index_range[item_key]:
            #     if len(item)>0:
            #         if item[0] < range_min_index:
            #             range_min_index = item[0]
            #         if item[1] > range_max_index:
            #             range_max_index = item[1]
            # if range_max_index ==0:
            #     reference_list_accept_index_range[item_key] = []
            # else:
            #     reference_list_accept_index_range[item_key] = [range_min_index,range_max_index]
                
        return reference_list_accept_index_range
    
    # 在回答里面添加引用标记
    def get_answer_with_reference_map(self,answer,reference_list,threshold=-1):
        print('输入的阈值',threshold)
        print('reference list 长度',len(reference_list))
        print('answer的长度',len(answer))
        if threshold > 0:
            self.threshold = threshold        
        answer_parts = self.split_answer(answer)
        print('answer_parts 的长度',sum([len(item) for item in answer_parts]))
        if len(answer) != sum([len(item) for item in answer_parts]):
            print('** 重大错误 **',answer)
        print(answer_parts)
        answer_parts_ref_indexs = []
        answer_parts_index = []
        for answer_part in answer_parts:
            if len(answer_parts_index)==0:
                answer_parts_index.append([0,len(answer_part)])
            else:
                answer_parts_index.append([answer_parts_index[-1][1],answer_parts_index[-1][1]+len(answer_part)])
            reference_indexs = self.find_sentence_reference(answer_part,reference_list=reference_list)
            answer_parts_ref_indexs.append(reference_indexs)
        print(answer_parts_index)
        
        reference_map = []
        # 合并句子引用
        for index in range(len(answer_parts)):
            reference_map.append({'index_range':answer_parts_index[index],'reference_list':answer_parts_ref_indexs[index]})
        
        reference_list_accept_index_range = self.find_part_index_from_reference(answer_parts,reference_list,reference_map)
        
        
        
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
        print(reference_map)
        self.threshold = self.init_threshold
        return reference_map,reference_list_accept_index_range