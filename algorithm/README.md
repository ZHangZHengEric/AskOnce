# AskOnce
全面而简单的知识库AI搜索系统-算法服务

standalone-docker-compose.yml 启动一个算法容器
[Dockerfile_build](Dockerfile_build) 是打包逻辑，会将对应服务的bin打包到bin目录下
[Dockerfile_run](Dockerfile_run) 是运行逻辑，依据bin目录和[supervisord.conf](supervisord.conf) 通过supervisord 启动

### supervisord.conf格式

```conf 
[supervisord]
nodaemon=true
env=WORKSPACE_DIR=/home/workspace/AskOnce/algorithm

[program:answer_by_documents]
command=/bin/bash -c "%(ENV_WORKSPACE_DIR)s/%(program_name)s.bin --tasktype %(program_name)s --worker_name %(program_name)s_%(process_num)s --api_key %(ENV_LLM_API_KEY)s --model_name %(ENV_LLM_MODEL_NAME)s --platform_api_url %(ENV_LLM_API_URL)s --search_url %(ENV_SEARCH_URL)s"
process_name=%(program_name)s_%(process_num)s
numprocs=1
autostart=true
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
...
```

一键部署规约

1. 新增的服务A必须在services/A/A.py ，并且taskType为A
2. 需要将新增的服务A添加到[Dockerfile_build](Dockerfile_build) 的构建脚本中
3. supervisord.conf中添加一个program:A,其他不用变更