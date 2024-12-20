FROM pytorch/pytorch:2.1.0-cuda12.1-cudnn8-devel as builder
ENV DEBIAN_FRONTEND=noninteractive

# 设置为中国国内源
RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list

WORKDIR /home/workspace

# 更新系统并安装必要的软件
RUN apt-get update -y --allow-unauthenticated && \
    apt-get install -y ttf-mscorefonts-installer && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 安装nuitka来打包
RUN pip install nuitka -i https://mirrors.aliyun.com/pypi/simple/ --no-cache-dir

# 设置 Python 工作路径
ENV PYTHONPATH /home/workspace:$PYTHONPATH

# 设置工作目录
WORKDIR /home/workspace/AskOnce/algorithm

# 复制当前代码到容器内并设置权限
COPY . .
RUN chmod 777 build.sh && ./build.sh

FROM pytorch/pytorch:2.1.0-cuda12.1-cudnn8-runtime
ENV DEBIAN_FRONTEND=noninteractive

# 设置为中国国内源
RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list
# 更新系统并安装必要的软件
RUN apt-get update -y --allow-unauthenticated && \
    apt-get install -y pandoc fonts-noto-cjk fontconfig libreoffice libreoffice-l10n-zh-cn \
    libreoffice-help-zh-cn git supervisor && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 安装服务所需的依赖包
RUN pip install --upgrade openai -i https://mirrors.aliyun.com/pypi/simple/ --no-cache-dir && \
    pip install scipy sentencepiece datasets rank_bm25 joblib sentence_transformers jionlp cn2an diffusers pymysql rouge safetensors --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple/ && \
    pip install weasyprint md2pdf opencc==1.1.6 pypandoc python-docx markdown python-pptx PyMuPDF tqdm unstructured matplotlib pandas python-magic pylibmagic pypdf opencc-python-reimplemented docx2txt langchain lxml lxml_html_clean trafilatura --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple/ && \
    pip install git+https://ghgo.xyz/https://github.com/mailgun/flanker@v0.9.14 --no-cache-dir && \
    pip install --upgrade "embedchain[dataloaders]" --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple/ && \
    pip install jieba elasticsearch Levenshtein --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple/

# 清理 apt 缓存
RUN apt-get purge -y --auto-remove && \
    apt-get clean

# 设置 Python 工作路径
ENV PYTHONPATH /home/workspace:$PYTHONPATH

WORKDIR /home/workspace/AskOnce/algorithm

# 复制时区等信息
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /etc/timezone /etc/timezone
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
# 复制构建结果
COPY --from=builder /home/workspace/AskOnce/algorithm/build .

COPY assets  /home/workspace/AskOnce

# 复制 supervisor 配置文件
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# 启动 supervisor 来管理多个进程
CMD ["supervisord", "-c", "/etc/supervisor/supervisord.conf"]

