import time
import random
'''
分布式唯一id: snowflake
'''
class Snowflake:
    def __init__(self, worker_id=1024, data_center_id=1024):
        # 机器标识 ID
        self.worker_id = worker_id
        # 数据中心 ID
        self.data_center_id = data_center_id
        # 计数序列号
        self.sequence = 0
        # 时间戳
        self.last_timestamp = -1
    
    def next_id(self):
        timestamp = int(time.time() * 1000)
        if timestamp < self.last_timestamp:
            raise Exception(f"Clock moved backwards. Refusing to generate id for {abs(timestamp - self.last_timestamp)} milliseconds")
        if timestamp == self.last_timestamp:
            self.sequence = (self.sequence + 1) & 4095
            if self.sequence == 0:
                timestamp = self.wait_for_next_millis(self.last_timestamp)
        else:
            self.sequence = 0
            self.last_timestamp = timestamp
        return ((timestamp - 1288834974657) << 22) | (self.data_center_id << 17) | (self.worker_id << 12) | self.sequence

