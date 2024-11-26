FROM pytorch/pytorch:2.1.0-cuda12.1-cudnn8-devel as builder
RUN rm -rf /home/workspace
WORKDIR /home/workspace
RUN apt-get update -y --allow-unauthenticated 
RUN apt install git -y
RUN pip install --upgrade openai -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install bitsandbytes==0.39.0 --upgrade -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install transformers --upgrade -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install peft --upgrade -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install accelerate --upgrade -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install deepspeed --upgrade -i https://mirrors.aliyun.com/pypi/simple/
# RUN pip install AutoGPTQ==0.3.2
RUN pip install scipy sentencepiece datasets rank_bm25 joblib sentence_transformers jionlp cn2an diffusers pymysql rouge safetensors -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install tritonclient[all] triton  gevent geventhttpclient -i https://mirrors.aliyun.com/pypi/simple/
RUN proxychains pip install flash-attn --no-build-isolation
# RUN DEBIAN_FRONTEND=noninteractive apt install libgtk-3-dev -y
# RUN DEBIAN_FRONTEND=noninteractive apt-get install libglib2.0-0 libglib2.0-dev python3-pip python3-cffi libffi-dev python3-brotli libpango-1.0-0 libpangoft2-1.0-0 -y
# # RUN ln -s /usr/lib/x86_64-linux-gnu/libgobject-2.0.so.0 /usr/local/lib/gobject-2.0-0.so.0
# RUN pip install pango -i https://mirrors.aliyun.com/pypi/simple/
# RUN pip install weasyprint --upgrade -i https://mirrors.aliyun.com/pypi/simple/
RUN conda install -c anaconda pango
RUN pip install weasyprint -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install md2pdf -i https://mirrors.aliyun.com/pypi/simple/
# RUN apt install -y --force-yes --no-install-recommends fonts-wqy-microhei ttf-wqy-zenhei
RUN apt update --fix-missing
RUN DEBIAN_FRONTEND=noninteractive apt-get install  ttf-mscorefonts-installer -y 
RUN apt install fonts-noto-cjk fontconfig -y
RUN apt install pandoc -y 
RUN pip install pypandoc -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install python-docx -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install matplotlib -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install lxml lxml_html_clean trafilatura   -i https://mirrors.aliyun.com/pypi/simple/
ENV PYTHONPATH /home/workspace:$PYTHONPATH
WORKDIR /home/workspace
ADD AskOnce /home/workspace/AskOnce
WORKDIR /home/workspace/AskOnce/algorithm/services
