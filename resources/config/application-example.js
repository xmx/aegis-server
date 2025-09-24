import config from 'aegis/server/config'
import console from 'console'
import time from 'time'

const cfg = {
    server: {
        addr: ":8443",
        read_timeout: 30 * time.second,
        read_header_timeout: time.parseDuration('10s'),
        static: {
            "/": "resources/static/root/",
            "/oas3": "resources/static/oas3/",
            "/play": "resources/static/play/"
        }
    },
    database: {
        uri: "mongodb+srv://username:password@mongo.example.com/dbname?retryWrites=true&w=majority&appName=dev"
    },
    logger: {
        level: "DEBUG",
        filename: "resources/log/application.jsonl"
    }
}

const level = cfg.logger.level
const terminal = cfg.logger.console
if (!terminal && level === 'DEBUG') {
    cfg.logger.console = true
    console.info(`日志级别是 ${level}，推测处于开发环境，已自动开启控制台输出日志`)
} else if (terminal && (level === 'ERROR' || level === 'WARN')) {
    cfg.logger.console = false
    console.info(`日志级别是 ${level}，推测处于生产环境，已自动开启控制台输出日志`)
}

config.set(cfg)
console.log('配置执行完毕')
