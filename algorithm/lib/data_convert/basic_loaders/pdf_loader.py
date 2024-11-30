import os
import re
import hashlib
from typing import Dict
from pypdf import PdfReader
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader

import fitz
import numpy as np
from tqdm import tqdm
class ReadSequenceByBBox:
    def __init__(self):
        pass
    def xy_cut(self,bboxes, direction="x"):
        result = []
        K = len(bboxes)
        indexes = range(K)
        if len(bboxes) <= 0:
            return result
        if direction == "x":
            # x first
            sorted_ids = sorted(indexes, key=lambda k: (bboxes[k][0], bboxes[k][1]))
            sorted_boxes = sorted(bboxes, key=lambda x: (x[0], x[1]))
            next_dir = "y"
        else:
            sorted_ids = sorted(indexes, key=lambda k: (bboxes[k][1], bboxes[k][0]))
            sorted_boxes = sorted(bboxes, key=lambda x: (x[1], x[0]))
            next_dir = "x"

        curr = 0
        np_bboxes = np.array(sorted_boxes)
        for idx in range(len(sorted_boxes)):
            if direction == "x":
                # a new seg path
                if idx != K - 1 and sorted_boxes[idx][2] < sorted_boxes[idx + 1][0]:
                    rel_res = self.xy_cut(sorted_boxes[curr:idx + 1], next_dir)
                    result += [sorted_ids[i + curr] for i in rel_res]
                    curr = idx + 1
            else:
                # a new seg path
                if idx != K - 1 and sorted_boxes[idx][3] < sorted_boxes[idx + 1][1]:
                    rel_res = self.xy_cut(sorted_boxes[curr:idx + 1], next_dir)
                    result += [sorted_ids[i + curr] for i in rel_res]
                    curr = idx + 1

        result += sorted_ids[curr:idx + 1]
        return result

    def augment_xy_cut(self,bboxes,
                    direction="x",
                    lambda_x=0.5,
                    lambda_y=0.5,
                    theta=5,
                    aug=False):
        if aug is True:
            for idx in range(len(bboxes)):
                vx = np.random.normal(loc=0, scale=1)
                vy = np.random.normal(loc=0, scale=1)
                if np.abs(vx) >= lambda_x:
                    bboxes[idx][0] += round(theta * vx)
                    bboxes[idx][2] += round(theta * vx)
                if np.abs(vy) >= lambda_y:
                    bboxes[idx][1] += round(theta * vy)
                    bboxes[idx][3] += round(theta * vy)
                bboxes[idx] = [max(0, i) for i in bboxes[idx]]
        res_idx = self.xy_cut(bboxes, direction=direction)
        res_bboxes = [bboxes[idx] for idx in res_idx]
        return res_idx, res_bboxes    



class PdfLoader(BaseLoader):
    def __init__(self, factory):
        super().__init__(factory)
        self.rsbbx = ReadSequenceByBBox()
    # def __init__(self, factory):
    #     super.__init__(factory)
    #     self.rsbbx = ReadSequenceByBBox()
    
    def load_data(self, url):
        result_text = ''
        self.cache_folder = os.environ.get('CONVERT_CACHE')
        detail_result  = []
        file_name = os.path.basename(url)
        if file_name.endswith('.pdf'):
            try:
                fitz_doc = fitz.open(url)
            except:
                print("文件损坏或者路径异常无法打开")
                return f"文件损坏或者路径异常无法打开", None
            
            # is_necrypted = fitz_doc.isEncrypted  # 是否pdf加密
            # if is_necrypted:
            #     print("该PDF文件已加密")
            #     return f"该PDF文件已加密", None
            file_name = file_name.replace('.pdf','')
            all_result = []
            start_page=0
            end_page = len(fitz_doc)
            print('pdf 一共',end_page,'页')
            print('开始页码',start_page,'结束页码',end_page)
            start_index = 0
            for page_index in tqdm(range(start_page,end_page)):
                page_width = int(fitz_doc[page_index].get_text('dict')['width'])
                page_height = int(fitz_doc[page_index].get_text('dict')['height'])
                page_dict = fitz_doc[page_index].get_text('dict')
                page_blocks = self.get_one_page_blocks(page_dict,page_index)
                block_ids,_  = self.rsbbx.augment_xy_cut([one_block[0] for one_block in page_blocks],direction="y")
                new_page_blocks = [ page_blocks[block_id] for block_id in block_ids ]
                page_blocks = new_page_blocks
                for one_block in page_blocks:
                    detail_result.append({'text':one_block[1],'text_index':len(detail_result)+1,'page_index':page_index+1,'bbox':one_block[0],'index_in_document':[start_index,start_index+len('，'+one_block[1])]})
                    start_index = start_index+len('，'+one_block[1])
                for one_block in page_blocks:
                    result_text += '，'+one_block[1]
                    # result_text += one_block[1]
        return result_text,detail_result
        
    def upscale_bbox(self,bbox,gaodu):
        bbox = [bbox[0]-gaodu,bbox[1]-gaodu,bbox[2]+gaodu,bbox[3]+gaodu]
        return bbox
    
    def isIntersection(self,bbox1,bbox2,gaodu=None):
        if gaodu:
            bbox1 = self.upscale_bbox(bbox1,gaodu)
            bbox2 = self.upscale_bbox(bbox2,gaodu)
        x01 = bbox1[0]
        y01 = bbox1[1]
        x02 = bbox1[2]
        y02 = bbox1[3]
        x11 = bbox2[0]
        y11 = bbox2[1]
        x12 = bbox2[2]
        y12 = bbox2[3]
        return (max(x01,x11)<=min(x02,x12)) and (max(y01,y11)<=min(y02,y12))
    
    def merge_blocks(self,blocks,is_upscale=False):
        new_blocks = []    
        new_block = blocks[0]
        for i in range(1,len(blocks)):
            # print('new_block[1]',new_block[1])
            # print('blocks[i][1]',blocks[i][1])
            # print(self.isIntersection(new_block[0], blocks[i][0]))
            if is_upscale:
                goadu = min((new_block[0][3]-new_block[0][1])/3,(blocks[i][0][3]-blocks[i][0][1])/3)
            else:
                goadu = None
            if self.isIntersection(new_block[0], blocks[i][0],goadu) and len(new_block[1])<128:
                new_block = [[min(new_block[0][0],blocks[i][0][0]),min(new_block[0][1],blocks[i][0][1]),max(new_block[0][2],blocks[i][0][2]),max(new_block[0][3],blocks[i][0][3]) ] , new_block[1]+blocks[i][1]]
            else:
                new_blocks.append(new_block)
                new_block = blocks[i]
        new_blocks.append(new_block)
        return new_blocks
    
    def get_one_page_blocks(self,page_block_dict,page_index):
        blocks = []
        images_num = 0
        for block in page_block_dict['blocks']:
            if block['type'] == 0:
                blocks_of_line = []
                for line in block['lines']:
                    for span in line['spans']:
                        if 'text' in span:
                            if len(span['text'].strip().replace(' ',''))>0:
                                blocks_of_line.append( [span['bbox'],span['text']] )
                                # blocks[-1][1] += span['text']
                        elif 'chars' in span:
                            for char_one in span['chars']:
                                if len(char_one['c'].strip().replace(' ',''))>0:
                                    blocks_of_line.append( [char_one['bbox'],char_one['c']] )
                                    # blocks[-1][1] += char_one['c']
                if len(blocks_of_line)>1:
                    blocks_of_line = self.merge_blocks(blocks_of_line,is_upscale=True)
                    page_blocks_old_len = len(blocks_of_line)
                    if page_blocks_old_len>0:
                        while True:
                            block_ids,_  = self.rsbbx.augment_xy_cut([one_block[0] for one_block in blocks_of_line],direction="y")
                            new_page_blocks = [ blocks_of_line[block_id] for block_id in block_ids ]
                            blocks_of_line = self.merge_blocks(new_page_blocks,is_upscale=True)
                            if len(blocks_of_line)==page_blocks_old_len or len(blocks_of_line)==0:
                                break
                            page_blocks_old_len = len(blocks_of_line)
                blocks.extend(blocks_of_line)     
            else:
                pass
        return blocks
    
        # reader = PdfReader(url)
        # data = []
        # all_content = []
        # meta_data = {"url": url, "file_size": f"{os.path.getsize(url)/ 1024 **2:.02f}MB", "file_type": url.split(".")[-1]}
        # number_of_pages = len(reader.pages)
        # if number_of_pages <= 0:
        #     raise ValueError("No data found")
        # for i in range(number_of_pages):
        #     page = reader.pages[i]
        #     page_meta_data = {}
        #     page_num = page.page_number
        #     content = page.extract_text()
        #     data.append(
        #         {
        #             "content": content,
        #             "page_num": page_num,
        #         }
        #     )
        #     all_content.append(content)
        # # 根据pdf内容生成id
        # all_content = "".join(all_content) 
        # doc_id = hashlib.sha256((all_content + url).encode()).hexdigest()
        # return all_content,{
        #     "id": doc_id,
        #     "meta_data": meta_data,
        #     "data": data,
        #     "text" : all_content
        # }
