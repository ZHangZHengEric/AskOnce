# AskOnce
全面而简单的知识库AI搜索系统-部署相关

环境变量设置如下：
```shell
# 后端配置目录
BACKEND_CONF_DIR= './conf'
BACKEND_MYSQL_ADDR='127.0.0.1:20000'
BACKEND_MYSQL_USER='root'
BACKEND_MYSQL_PASSWORD='123456'
BACKEND_REDIS_ADDR='127.0.0.1:6379'
BACKEND_REDIS_PASSWORD='123456'
BACKEND_MINIO_ADDR='http://127.0.0.1:20002'
BACKEND_MINIO_AK=''
BACKEND_MINIO_SK=''
BACKEND_ES_ADDR='http://127.0.0.1:20004'
############################################
# minio数据存储目录
MINIO_DATA_DIR="/data/minio"
# minio root密码
MINIO_ROOT_PASSWORD='123456'

#############################################
# msyql配置目录
MYSQL_CONF_DIR="./conf/mysql"
# msyql数据存储目录
MYSQL_DATA_DIR="/data/mysql"
# mysql root密码
MYSQL_ROOT_PASSWORD="123456"

#############################################
# redis数据存储目录
REDIS_DATA_DIR="/data/redis"
# redis配置目录
REDIS_CONF_DIR="./conf/redis"
# redis root密码
REDIS_ROOT_PASSWORD='123456'

#############################################
# elasticsearch数据存储目录
ES_DATA_DIR="/data/elasticsearch"
# elasticsearch配置目录
ES_CONF_DIR="./conf/elasticsearch"
# elasticsearch用户名
ES_USERNAME="elastic"
# elasticsearch密码
ES_PASSWORD="123456"
# es地址
ES_ADDR="http://127.0.0.1:20004"

##############################################

# jobd地址
JOBD_ADDR="http://127.0.0.1:20033"

# LLM API INFO
LLM_API_KEY="sk-0d772fccdde0421088e31668e862eee9"
LLM_MODEL_NAME="deepseek-chat"
LLM_API_URL="https://api.deepseek.com/v1"

# CONVERT CACHE FOR SAVING TEMP DATA
CONVERT_CACHE="/srv/temp_convert_cache"

```
## 安装说明
1. elasticsearch安装：设置数据目录权限：
mkdir -p ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3
chmod g+rwx ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3
chgrp 0 ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3

2. elasticsearch安装：配置max_map_count
   sysctl -w vm.max_map_count=262144
