> 要使用邮件服务，需要在 `.env` 文件中配置 SMTP 服务

### 发送登陆密码重置邮件

[POST] /v1/email/send/password/reset

发送登陆密码重置邮件

| 参数  | 类型     | 说明     | 必选 |
| ----- | -------- | -------- | ---- |
| email | `string` | 邮箱地址 | \*   |

### 发送注册邮件

[POST] /v1/email/send/register

发送注册帐号的邮件

| 参数         | 类型     | 说明                                       | 必选 |
| ------------ | -------- | ------------------------------------------ | ---- |
| email        | `string` | 邮箱地址                                   | \*   |
| redirect_url | `string` | 邮件内容的跳转链接，用户点击之后跳转的链接 | \*   |

`redirect_url` 会带有 `?code=xxx&email=xxx` 传给前端.

前端可获取这两个参数，完成注册
