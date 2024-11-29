# AskOnce
全面而简单的知识库AI搜索系统-部署相关

环境变量设置如下：
```shell
# minio数据存储目录
MINIO_DATA_DIR="./data/minio"

# msyql数据存储目录
MYSQL_DATA_DIR="./data/mysql"
MYSQL_CONF_DIR="./conf/mysql"
# redis数据存储目录
REDIS_DATA_DIR="./data/redis"
REDIS_CONF_DIR="./conf/redis"
# elasticsearch数据存储目录
ES_DATA_DIR="./data/elasticsearch"
ES_CONF_DIR="./conf/elasticsearch"

```
## 安装说明
1. elasticsearch安装：设置数据目录权限：
mkdir -p ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3
chmod g+rwx ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3
chgrp 0 ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3

2. elasticsearch安装：配置max_map_count
   sysctl -w vm.max_map_count=262144
