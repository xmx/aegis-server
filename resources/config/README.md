配置文件支持三种格式：js/json/jsonc。

JSONC（即 **JSON** with **C**omments）是对标准 JSON 的一种扩展，允许在 JSON 文档中加入注释。

## js 方式

```js
import config from 'aegis/server/config'

// 配置参数通过 config 操作
// config.get() 获取配置参数
// config.set(cfg) 设置配置参数
```
