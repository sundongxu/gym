```mermaid
sequenceDiagram
    participant C as 小程序端
    participant S as 后端服务
    participant W as 微信服务器
    participant D as 数据库

    C->>C: 1. 调用 wx.login() 获取 code
    C->>S: 2. 请求登录 POST /api/v1/login/wx-mini<br/>Body: {code: "xxx"}
    S->>W: 3. 请求 code2session 接口<br/>携带 appid + secret + code
    W-->>S: 4. 返回 openid + session_key
    
    S->>D: 5. 查询 user_auths 表<br/>WHERE login_type='wx_mini_app' AND auth_key={openid}
    
    alt 首次登录-用户不存在
        S->>D: 6.1 创建新用户记录 users 表
        S->>D: 6.2 创建认证记录 user_auths 表
    else 非首次登录
        S->>D: 6.3 更新 auth_secret 和 last_login
    end
    
    S->>S: 7. 生成 JWT token
    S-->>C: 8. 返回登录结果<br/>{"token": "xxx", "user": {...}}