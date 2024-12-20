FROM pytorch/pytorch:2.1.0-cuda12.1-cudnn8-devel as builder
WORKDIR /home/workspace
# 更新系统并安装必要的软件
RUN apt-get update -y --allow-unauthenticated
RUN apt install git -y
RUN apt update --fix-missing
RUN DEBIAN_FRONTEND=noninteractive apt-get install ttf-mscorefonts-installer -y
# 安装pyinstaller来打包
RUN pip install pyinstaller -i https://mirrors.aliyun.com/pypi/simple/
# 设置 Python 工作路径
ENV PYTHONPATH /home/workspace:$PYTHONPATH
#设置工作目录
WORKDIR /home/workspace/AskOnce/algorithm
COPY . .
# 安装服务需要的所有依赖
RUN pip install --upgrade openai -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install scipy sentencepiece datasets rank_bm25 joblib sentence_transformers jionlp cn2an diffusers pymysql rouge safetensors -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install weasyprint -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install md2pdf -i https://mirrors.aliyun.com/pypi/simple/
RUN apt install fonts-noto-cjk fontconfig -y
RUN apt install libreoffice libreoffice-l10n-zh-cn libreoffice-help-zh-cn -y --fix-missing
RUN pip install opencc==1.1.6 -i https://mirrors.aliyun.com/pypi/simple/
RUN apt install pandoc -y
RUN pip install pypandoc -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install python-docx markdown python-pptx PyMuPDF tqdm unstructured -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install matplotlib pandas python-magic pylibmagic pypdf opencc-python-reimplemented docx2txt langchain -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install lxml lxml_html_clean trafilatura   -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install git+https://ghgo.xyz/https://github.com/mailgun/flanker@v0.9.14
RUN pip install --upgrade "embedchain[dataloaders]"
RUN pip install jieba -i https://mirrors.aliyun.com/pypi/simple/
RUN pip install elasticsearch Levenshtein -i https://mirrors.aliyun.com/pypi/simple/
