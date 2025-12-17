**登陆成功后 后端POST方法返回**

```json
{
    "code": 200
    "message": 信息
    "data": 
    {
    	"token_class":
    	{
    		"token": token
    		"refresh_token": 用于刷新token
    		"expires_in": token到期时间（以秒为单位）
		}
    	"user_info": 
    	{
    		"uid": 微信号
    		"name": 昵称
		}
	}
}
```

**refresh_token前端接口**

```http
POST /auth/refresh_token
Authorization: Bearer <refresh_token>
```

**refresh_token 成功后 后端POST方法返回**

```json
{
    "code": 
    "message": 信息
    "data":
    {
    	"token": 刷新之后的token
    	"refresh_token": 新的用于刷新的token
    	"expires_in": token到期时间（秒）
	}
}
```

**refresh_token失败，返回**

```json
{
    "code": 状态码
    "message": 信息
}
```

**refresh_token 状态码**

```html
401: token格式错误或失效
400: 请求参数错误
409: 请求有冲突
```



























