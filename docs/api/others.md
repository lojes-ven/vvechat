**注册**

前端接口：

```http
POST /api/register
Content-Type: application/json

{
    "name": "用户名",
    "password": "密码",
    "phone_number": "手机号"
}
```

成功后端返回：

```json
{
    "code": 201,
    "message": "注册成功"
}
```

失败后端返回：

```json
{
    "code": 400,
    "message": "手机号已存在" // 或其他错误信息
}
```

**微信号登录**

前端接口：

```http
POST /api/login/uid
Content-Type: application/json

{
    "uid": "微信号",
    "password": "密码"
}
```

成功后端返回：

```json
{
    "code": 200,
    "message": "登陆成功",
    "data": {
        "user_info": {
            "name": "用户名",
            "uid": "微信号"
        },
        "token_class": {
            "token": "JWT token",
            "refresh_token": "刷新token",
            "expires_in": 3600
        }
    }
}
```

失败后端返回：

```json
{
    "code": 400,
    "message": "错误信息"
}
```

**手机号登录**

前端接口：

```http
POST /api/login/phone_number
Content-Type: application/json

{
    "phone_number": "手机号",
    "password": "密码"
}
```

成功后端返回：同微信号登录

失败后端返回：同微信号登录

**刷新Token**

前端接口：

```http
POST /api/auth/refresh_token
Authorization: Bearer <token>
```

成功后端返回：

```json
{
    "code": 201,
    "message": "success",
    "data": {
        "token": "新JWT token",
        "refresh_token": "新刷新token",
        "expires_in": 3600
    }
}
```

失败后端返回：

```json
{
    "code": 500,
    "message": "token出现问题..."
}
```

**修改微信号**

前端接口：

```http
POST /api/auth/me/uid
Authorization: Bearer <token>
Content-Type: application/json

{
    "uid": "新微信号"
}
```

成功后端返回：

```json
{
    "code": 201,
    "message": "success"
}
```

失败后端返回：

```json
{
    "code": 400,
    "message": "微信号重复"
}
```

**修改昵称**

前端接口：

```http
POST /api/auth/me/name
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "新昵称"
}
```

成功后端返回：

```json
{
    "code": 200, // 或201，视具体实现而定，通常成功即可
    "message": "success"
}
```

**修改密码**

前端接口：

```http
POST /api/auth/me/password
Authorization: Bearer <token>
Content-Type: application/json

{
    "prev_password": "旧密码",
    "new_password": "新密码"
}
```

成功后端返回：

```json
{
    "code": 200,
    "message": "success"
}
```

**查看好友信息**

前端接口：

```http
GET /api/auth/info/friends/id/{id}
Authorization: Bearer <token>
```

成功后端返回：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "好友ID",
        "friend_remark": "备注",
        "name": "昵称",
        "uid": "微信号"
    }
}
```

失败后端返回：

```json
{
    "code": 400,
    "message": "找不到好友"
}
```

**查看陌生人信息**

前端接口：

```http
GET /api/auth/info/strangers/id/{id}
Authorization: Bearer <token>
```

成功后端返回：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "陌生人ID",
        "name": "昵称"
    }
}
```

失败后端返回：

```json
{
    "code": 400,
    "message": "找不到此人"
}
```

**查看好友信息通过Uid**

前端接口：

```http
GET /api/auth/info/friends/uid/{uid}
Authorization: Bearer <token>
```

成功后端返回：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "好友ID",
        "friend_remark": "备注",
        "name": "昵称",
        "uid": "微信号"
    }
}
```

失败后端返回：

```json
{
    "code": 400,
    "message": "找不到好友"
}
```

**查看陌生人信息通过Uid**

前端接口：

```http
GET /api/auth/info/strangers/uid/{uid}
Authorization: Bearer <token>
```

成功后端返回：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "陌生人ID",
        "name": "昵称"
    }
}
```

失败后端返回：

```json
{
    "code": 400,
    "message": "找不到此人"
}
```

