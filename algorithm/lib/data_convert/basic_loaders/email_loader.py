import os
import re
import uuid
import json
import base64
import email
import hashlib
import logging
import chardet
import traceback
from tqdm import tqdm
from glob import glob
from email import policy
from opencc import OpenCC
from textwrap import dedent
from datetime import datetime
from bs4 import BeautifulSoup
from email.parser import Parser
from pathlib import Path
from flanker.addresslib import address
from functools import reduce
from dataclasses import dataclass
from email.utils import parsedate_to_datetime
from email.parser import BytesParser
from email.feedparser import headerRE
from email.iterators import _structure
from typing import Dict, Optional, List
from email.header import decode_header
from AskOnce.algorithm.lib.data_convert.utils.dist_id import Snowflake
from AskOnce.algorithm.lib.data_convert.utils.common_utils import clean_string
from AskOnce.algorithm.lib.data_convert.basic_loaders.base_loader import BaseLoader
from AskOnce.algorithm.lib.data_convert.utils.json_serializable import register_deserializable, JSONSerializable
from AskOnce.algorithm.lib.data_convert.utils.dist_id import Snowflake
# 本地调试, push前备注掉

try:
    CONVERT_CACHE = os.environ.get("CONVERT_CACHE")
    if CONVERT_CACHE is None or len(CONVERT_CACHE) <=0:
        raise Exception("CONVERT_CACHE 不存在")
except:
    raise Exception("CONVERT_CACHE 不存在")
    
dist_id = Snowflake(worker_id=2024, data_center_id=2024)
convert = OpenCC("tw2s").convert

def clean_body_text_placeholder(body_text:str):
    # 去掉一些邮件图片解析的一些乱码占位符
    body_text = html2txt(body_text)
    body_text = re.sub(r'<mailto:.*?>', '', body_text)
    body_text = re.sub(r'\[cid:.*?\]','', body_text)
    body_text = re.sub(r'\bmage.*?','', body_text)       # 去
    body_text = re.sub(r'=[=~]+(?![=~])','', body_text)  # 去 =~=~=~=~=~=~=~=~=~=~=~=~=~=
    body_text = re.sub(r"---{1,}",'', body_text)
    body_text = re.sub(r'___{1,}','', body_text)
    body_text = re.sub(r'====={1,}','', body_text)
    
    return body_text

def clean_whitespace(text):
     # 多个连续空格 变为 1个空格, 不能全去掉，会有英/韩/日文
    text = re.sub(r"\s+"," ",text)
    return text


def html2txt(html_content):
    soup = BeautifulSoup(html_content, 'html.parser')
    text = soup.get_text(strip=True)
    return text

class EmailLoader(BaseLoader):
    
    def decode_header(self, header_item):
        '''
        from/ to有形式：
        "Tuan Hung Ngo" <TUANHUNG@medigenvac.com>  西语名
        "MeeiYun林美雲" <MeeiYun@medigenvac.com>   中西 /中
        vtd@nihe.org.vn                       没有名字
        subject 有二种形式：
        Subject: Response from the discussion on May 15
        '''
        if header_item == None: return ""
        items = decode_header(header_item)  # List[Dict(内容, 编码)]
        decoded_header = []
        for item in items:
            item_value = item[0] if isinstance(item[0], str) == True else item[0].decode(item[1] or "utf-8")
            decoded_header.append(item_value)
        return decoded_header
    
    
    @staticmethod
    def _process_html_content(html_content) -> str:
        '''邮件html正文 解析'''
        content = BeautifulSoup(html_content, "html.parser").get_text()
        # content = clean_string(content)
        return dedent(
            f"""
            Content: {content}
        """
        )         
    
    @staticmethod
    def _decode_payload(part):
        charset = part.get_content_charset() or "utf-8"
        if charset == 'gb2312':
            charset = 'GB18030'
        try:
            if  part['Content-Transfer-Encoding'] in ['base64','quoted-printable']:
                if part['Content-Transfer-Encoding'] in ['base64']:
                    payload = part.get_payload(decode=False)
                    # 检查payload 中是否是连续的base64 编码，不是清楚掉其他的内容
                    payload_list = payload.split('\n')
                    for index, item in enumerate(payload_list):
                        if len(item)==0:
                            payload_list = payload_list[:index]
                            break
                    return base64.b64decode('\n'.join(payload_list)).decode(charset,errors='ignore').encode('utf-8').decode('utf-8')
                return part.get_payload(decode=True).decode(charset,errors='ignore').encode('utf-8').decode('utf-8')
            return part.get_payload(decode=False).encode('utf-8').decode('utf-8')
        except Exception as e:
            print(traceback.print_exc())
            print('error')
            return part.get_payload(decode=True).decode(charset, errors="replace").encode('utf-8').decode('utf-8')
    
    def get_attachment(self, mime_msg):
        '''邮件附件解析'''
        attachment = {
            "files":[],
            "images":[],
        }
        for attach in mime_msg.iter_attachments():
            ctype = attach.get_content_type()
            charset = attach.get_charset()
            file_name = attach.get_filename() #如果是附件，这里就会取出附件的文件名
            if file_name is None:
                file_name = "attach.pdf"  # 需要猜测attach_data的文件类型
            local_file_path = os.path.join(os.environ.get("CONVERT_CACHE"), file_name)
            attach_data = attach.get_payload(decode=True)
            if attach_data and ctype != "message/rfc822":
                attach_content = ""
                file_name_ = file_name[-8:]
                local_file_path = os.path.join(CONVERT_CACHE,str(dist_id.next_id()) + file_name_)
                with open(local_file_path, 'wb') as f:
                    f.write(attach_data)
                    # print(f'附件文档已保存: {local_file_path}')
                    if Path(local_file_path).suffix in ['.pdf','.docx', '.txt']:
                        # attach_content, _ = self._get_payload_content(local_file_path)
                        try:
                            attach_content, _ = self.factory.create(local_file_path)
                            # 在UTF-8编码中，代理字符是一些特殊的字符，用于在UTF-16编码中表示那些在基本多语言平面（BMP）之外的字符。这些字符在UTF-8编码中是不允许出现的，因此在尝试将包含这些代理字符的文本编码为UTF-8时会出现错误。
                            attach_content = attach_content.encode('utf-8', 'ignore').decode('utf-8','ignore')
                            # print('attach_content',attach_content)
                            attach_content = convert(attach_content)
                            # attach_content = clean_whitespace(attach_content)
                        except Exception as e:
                            traceback.print_exc()
                            attach_content, _ = "", None
            else:
                attach_content = ""
                print("附件内容为None或者附件是邮件处理不了") # CONVERT_CACHE/ssxxx.pdf
            attachment['files'].append({"file_name":local_file_path, "file_content":attach_content, "attach_name": file_name})
            if ctype.startswith("image/"):
                file_name = attach.get_filename()
                if file_name is None:
                    file_name = '图片.png' 
                # 解码邮件头部的非ASCII字符
                # file_name = decode_header(file_name)[0][0]
                # if isinstance(file_name, bytes):
                #     file_name = file_name.decode()
                # 保存附件到本地
                local_file_path = os.path.join(CONVERT_CACHE, str(dist_id.next_id()) + file_name)
                with open(local_file_path, 'wb') as f:
                    f.write(attach.get_payload(decode=True))
                # print(f'图片已保存: {local_file_path}')
                attachment["images"].append(local_file_path)
        return attachment
    
    def get_body(self, mime_msg):
        '''邮件正文解析'''
        body = {
            "html_content":"",
            "html2txt_content":"",
            "text_content":"",
        }
        if mime_msg.is_multipart():
            for part in mime_msg.walk():
                ctype = part.get_content_type()
                file_name = part.get_filename()
                if ctype == "text/html":
                    html_content = self._decode_payload(part)
                    html2txt_content = html2txt(html_content)
                    body['html_content'] += convert(html_content)
                    body['html2txt_content'] += html2txt_content
                if ctype == "text/plain":
                    text_content =  self._decode_payload(part)
                    text_content = clean_body_text_placeholder(text_content)
                    # text_content = clean_whitespace(text_content)
                    body['text_content'] += convert(text_content)
        else:
            text_content = self._decode_payload(mime_msg)
            text_content = clean_body_text_placeholder(text_content)
            body['text_content'] = convert(text_content)
        return body
  

    def load_data(self, file_path):
        encodeing_type = ['utf-8','big5','gbk','gb2312']
        for item in  encodeing_type:
            try:
                f = open(file_path, 'r',encoding=item) # errors='ignore'
                eml_content = f.read()
            except:
                f = None
                print(item,'not work')
                continue
            break
        if f is None:
            f = open(file_path, 'r',encoding='utf-8',errors='ignore') # errors='ignore'
            eml_content = f.read()
        # raw_email = open(file_path, 'rb').read()
        # detected_encoding = chardet.detect(raw_email)['encoding']
        # print(detected_encoding)
        # eml_content = raw_email.decode(detected_encoding)
        # print(eml_content)
        eml_content_lines = eml_content.split('\n')
        start_index=0
        end_index = len(eml_content_lines)
        for index, line in enumerate(eml_content_lines):
            if headerRE.match(line):
                start_index = index
                break

        eml_content = '\n'.join(eml_content_lines[start_index:end_index])
        # smtp_pattern = re.compile(r'MAIL FROM:<[^>]+>.*?BDAT \d+ LAST\n', re.DOTALL)
        # smtp_commands = smtp_pattern.findall(eml_content)

        # 移除SMTP命令后，剩下的就是邮件内容
        # email_content = smtp_pattern.sub('', eml_content)
        # print(eml_content)
        # mime_msg = BytesParser(policy=policy.default).parsebytes(mail_data)
        mime_msg = email.message_from_string(eml_content, policy=policy.default)
        
        # print('邮件的key',mime_msg.keys())
        # 邮件结构
        # print(_structure(mime_msg))
        # dfs 邮件结构, 搞清楚 content_type, 编码, 
        # for part in mime_msg.walk():
        #     # print(part.items())
        #     print(part.get_content_type())
        #     print(part.get_content_charset())
        try:
            sender = convert(mime_msg.get("From", ''))
            sender_obj = address.parse(sender)
            try:
                sender_name = sender_obj.display_name
            except:
                sender_name = ""
            try:
                sender_address = sender_obj.ace_address
            except:
                sender_address = ""
        except:
            sender = ""
            sender_name = ""
            sender_address = ""
            
        try:
            receiver= convert(mime_msg.get("To",''))
        except:
            receiver = ''
        receiver_list = receiver.split(",")
        receiver_obj_list = address.parse_list(receiver_list)
        receiver_names = []
        receiver_addresses = []
        for item in receiver_obj_list:
            try:
                receiver_names.append(item.display_name)
                receiver_addresses.append(item.ace_address)
            except:
                continue
        subject= convert(mime_msg.get("Subject",''))
        try:
            date_str = mime_msg.get("Date", '')
            date_obj = parsedate_to_datetime(date_str)
            date_str = datetime.strftime(date_obj, "%Y-%m-%d %H:%M:%S")
            date = convert(date_str)
        except:
            date = ''
        cc = convert(mime_msg.get("CC", ''))
        body = self.get_body(mime_msg) # 保留换行格式
        if body['text_content'] is None or len(body['text_content']) <= 0:
            if body['html2txt_content'] is not None and len(body['html2txt_content']) > 0:
                body['text_content'] = body['html2txt_content']
        body.pop("html2txt_content")
        attachment = self.get_attachment(mime_msg)
        strings = ["发件人：" + sender, "收件人：" + receiver, "抄送：" + cc , "邮件主题：" + subject , "收发讯息时间：" + date , "邮件正文：" + body['text_content']]
        # text中去掉发件抄送人
        # strings = ["邮件主题：" + subject , "收发讯息时间：" + date , "邮件正文：" + body['text_content']]
        # strings = ["邮件主题：" + subject ,"邮件正文：" + body['text_content']]
        text = reduce(lambda x, y: x + "\n" + y, strings) 
        # text = '\n'.join(strings)
        text_detail = {
            "sender": sender,
            "receiver": receiver,
            'sender_name': sender_name,
            'sender_address': sender_address,
            'receiver_names': receiver_names,
            'receiver_addresses': receiver_addresses,
            "date": date,
            "subject": subject,
            "cc": cc,
            "body": body,
            "attachment":attachment
        }
        return text, text_detail