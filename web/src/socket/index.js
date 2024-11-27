/**
 * create by zhangxiang on 2022-03-05 19:58
 * 类注释：
 * 备注：
 */
import {useSocketStore} from '@/store'

const socketStore = useSocketStore()

export default class SocketService {
    /**
     * 单例
     */
    static instance = null;

    static get Instance() {
        if (!this.instance) {
            this.instance = new SocketService()
        }
        return this.instance
    }

    // 和服务端连接的socket对象
    ws = null

    // 标识是否连接成功
    connected = false

    // 记录重试的次数
    sendRetryCount = 0

    // 重新连接尝试的次数
    connectRetryCount = 0


    //  定义连接服务器的方法
    connect(url) {
        // 连接服务器
        if (!window.WebSocket) {
            return console.log('您的浏览器不支持WebSocket')
        }
        this.ws = new WebSocket(url)
        // 连接成功的事件
        this.ws.onopen = () => {
            console.log('连接服务端成功')
            this.connected = true
            // 重置重新连接的次数
            this.connectRetryCount = 0
            this.heartReset()
        }
        // 1.连接服务端失败
        // 2.当连接成功之后, 服务器关闭的情况(连接失败重连)
        this.ws.onclose = (e) => {
            console.log('连接服务端失败', e)
            console.log('websocket 断开: ' + e.code + ' ' + e.reason + ' ' + e.wasClean)
            this.connected = false
            this.connectRetryCount++
            setTimeout(() => {
                if (e.code !== 1005) {
                    this.connect(url)
                }
            }, 3 * 1000 * this.connectRetryCount)
        }
        // 得到服务端发送过来的数据
        this.ws.onmessage = msg => {
            this.heartReset()
            msg = JSON.parse(msg.data)
            const msgType = msg.msg_type
            const content = JSON.parse(msg.content)
            const newData = {...content, msgType}
            switch (msgType) {
                case 'ws_chat':
                    socketStore.setWsChat(newData)
                    break
                case 'ws_task_num':
                    socketStore.setWsTaskNum(newData)
                    break
                default:
                    console.log('不匹配')
            }
        }
    }

    // 发送数据的方法
    send(data) {
        // 判断此时此刻有没有连接成功
        if (this.connected) {
            this.sendRetryCount = 0
            this.ws.send(data)
        } else {
            this.sendRetryCount++
            setTimeout(() => {
                this.send(data)
            }, this.sendRetryCount * 500)
        }
    }

    close() {
        if (this.connected) {
            console.log('关闭连接')
            this.ws.close()
        }
    }

    // 服务器和客户端收到的信息内容如果如下 就识别为心跳信息 不要做业务处理
    checkMsg = 'heartbeat';
    // 延时发送消息对象（启动心跳新建这个对象，收到消息后重置对象）
    timeoutObj = null

    // 一段时间后发送心跳包
    heartStart() {
        this.timeoutObj = setInterval(() => {
            if (this.connected) this.ws.send(this.checkMsg)
        }, 25 * 1000)
    }

    // 接收到服务器的信息之后要重置心跳发送的方法
    heartReset() {
        clearInterval(this.timeoutObj)
        this.heartStart()
    }
}

