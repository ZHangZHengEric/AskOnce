# 第一阶段：构建和安装依赖
FROM pytorch/pytorch:2.1.0-cuda12.1-cudnn8-devel as LLM_builder

# 设置工作目录
WORKDIR /home/workspace

# 更新系统并安装必要工具和依赖
RUN apt-get update -y --allow-unauthenticated && \
    apt-get install -y --no-install-recommends \
        git \
        fonts-noto-cjk \
        fontconfig \
        libreoffice libreoffice-l10n-zh-cn libreoffice-help-zh-cn \
        ttf-mscorefonts-installer \
        pandoc && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 设置 Python 镜像源，并安装所有 Python 依赖
RUN pip install --upgrade pip -i https://mirrors.aliyun.com/pypi/simple/ && \
    pip install --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple/ \
        openai \
        scipy \
        sentencepiece \
        datasets \
        rank_bm25 \
        joblib \
        sentence_transformers \
        jionlp \
        cn2an \
        diffusers \
        pymysql \
        rouge \
        safetensors \
        weasyprint \
        md2pdf \
        opencc==1.1.6 \
        pypandoc \
        python-docx \
        markdown \
        python-pptx \
        PyMuPDF \
        tqdm \
        unstructured \
        matplotlib \
        pandas \
        python-magic \
        pylibmagic \
        pypdf \
        opencc-python-reimplemented \
        docx2txt \
        langchain \
        lxml \
        lxml_html_clean \
        trafilatura \
        jieba \
        elasticsearch \
        Levenshtein && \
    pip install --no-cache-dir git+https://github.com/mailgun/flanker@v0.9.14 && \
    pip install --no-cache-dir --upgrade "embedchain[dataloaders]"

# 第二阶段：构建运行镜像
FROM pytorch/pytorch:2.1.0-cuda12.1-cudnn8-runtime as runtime

# 设置工作目录
WORKDIR /home/workspace

# 更新系统并安装必要的工具
RUN apt-get update -y --allow-unauthenticated && \
    apt-get install -y --no-install-recommends \
    fonts-noto-cjk \
    fontconfig \
    supervisor \
    wget \
    gnupg2 \
    lsb-release \
    curl \
    unzip \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# 安装 MySQL 8.0
RUN wget https://dev.mysql.com/get/mysql-apt-config_0.8.17-1_all.deb && \
    dpkg -i mysql-apt-config_0.8.17-1_all.deb && \
    apt-get update && \
    apt-get install -y mysql-server=8.0* && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 安装 Redis
RUN curl -fsSL https://packages.redis.io/gpg | gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/redis.list && \
    apt-get update && \
    apt-get install -y redis && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 复制依赖和项目文件
COPY --from=LLM_builder /home/workspace /home/workspace

# 设置 Python 环境变量
ENV PYTHONPATH /home/workspace:$PYTHONPATH

# 设置最终工作目录
WORKDIR /home/workspace/AskOnce/algorithm/services

# 复制 supervisor 配置文件
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# 启动 supervisor 来管理多个进程
CMD ["supervisord", "-c", "/etc/supervisor/supervisord.conf"]

