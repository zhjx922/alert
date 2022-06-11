# LogAlert

重量只有1克的轻量级日志监控告警程序

## 功能

* 日志文件配置支持glob
* 监控关键字配置
* 告警CURL配置

## 运行

```yaml
go run main.go -c config_demo.yaml
```

## 配置说明

使用yaml格式

```yaml
inputs:
    # 项目名称，没啥用
  - name: project-000
    # 监控的文件，支持glob
    paths:
      - /data/log/*.log
      - /data/log1/log1.log
    # 监控内容，包含内容即报警
    include_lines: ['error', 'warning']
    # 另一个项目配置
  - name: project-001
    paths:
      - /var/log/*.log
    include_lines: ['success']
output.http:
  method: POST
  # 这里是企微机器人的地址
  url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=*
  # Header头
  headers:
    - Content-Type application/json;charset=UTF-8
  # 留着扩展用
  format: json
  # 请求内容(%{content}会替换为日志告警行的内容)
  body: >
    {
      "msgtype": "markdown",
      "markdown": {
        "content": "DIY报警内容\n<font color=\"warning\">%{content}</font>"
      }
    }
```