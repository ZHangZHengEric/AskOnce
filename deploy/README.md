# **AskOnce**
> **全面而简单的知识库 AI 搜索系统 - 部署相关**

---

## **环境变量设置**

以下是系统所需的环境变量配置，确保在部署前正确设置：

```shell
# ========================= 后端配置 =========================
# 后端配置目录
BACKEND_CONF_DIR='./conf'
BACKEND_MYSQL_ADDR='127.0.0.1:20000'
BACKEND_MYSQL_USER='root'
BACKEND_MYSQL_PASSWORD='123456'
BACKEND_REDIS_ADDR='127.0.0.1:6379'
BACKEND_REDIS_PASSWORD='123456'
BACKEND_MINIO_ADDR='http://127.0.0.1:20002'
BACKEND_MINIO_AK=''
BACKEND_MINIO_SK=''
BACKEND_ES_ADDR='http://127.0.0.1:20004'

# ========================= MinIO 配置 =========================
# MinIO 数据存储目录
MINIO_DATA_DIR="/data/minio"
# MinIO Root 密码
MINIO_ROOT_PASSWORD='123456'

# ========================= MySQL 配置 =========================
# MySQL 配置目录
MYSQL_CONF_DIR="./conf/mysql"
# MySQL 数据存储目录
MYSQL_DATA_DIR="/data/mysql"
# MySQL Root 密码
MYSQL_ROOT_PASSWORD="123456"

# ========================= Redis 配置 =========================
# Redis 数据存储目录
REDIS_DATA_DIR="/data/redis"
# Redis 配置目录
REDIS_CONF_DIR="./conf/redis"
# Redis Root 密码
REDIS_ROOT_PASSWORD='123456'

# ========================= Elasticsearch 配置 =========================
# Elasticsearch 数据存储目录
ES_DATA_DIR="/data/elasticsearch"
# Elasticsearch 配置目录
ES_CONF_DIR="./conf/elasticsearch"
# Elasticsearch 用户名与密码
ES_USERNAME="elastic"
ES_PASSWORD="123456"
# Elasticsearch 地址
ES_ADDR="http://127.0.0.1:20004"

# ========================= 其他配置 =========================
# Jobd 地址
JOBD_ADDR="http://127.0.0.1:20033"

# LLM API 信息
LLM_API_KEY="sk-**"
LLM_MODEL_NAME="deepseek-chat"
LLM_API_URL="https://api.deepseek.com/v1"

# 转换缓存目录 (用于保存临时数据)
CONVERT_CACHE="/srv/temp_convert_cache"

```
---

## **安装说明**

### **1. Elasticsearch 安装**

#### **设置数据目录权限**
```bash
mkdir -p ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3
chmod g+rwx ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3
chgrp 0 ${ES_DATA_DIR}/node1 ${ES_DATA_DIR}/node2 ${ES_DATA_DIR}/node3
```
#### **配置max_map_count**
```bash
sysctl -w vm.max_map_count=262144
```
---

## **启动说明**

整体 `docker-compose` 文档根据功能分为以下几类 **profiles**：

| **Profile** | **描述**                 |
|-------------|--------------------------|
| `all`       | 启动全部服务              |
| `nginx`     | 前端服务                 |
| `backend`   | 后端业务服务             |
| `db`        | 数据库等持久化服务        |
| `llm`       | 算法服务                 |

### **启动命令示例**

#### **1. 启动全部服务**
```bash
docker compose --profile all up -d
```
#### **1. 启动对应服务**
```bash
docker compose --profile llm up -d
```

以上内容旨在帮助快速部署和启动 AskOnce 系统。如有问题，请参考官方文档或联系技术支持。





