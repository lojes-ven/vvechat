**登陆成功后后端返回**

```json
{
    "code": http状态码
    "message": 信息
    "data": 
    {
    	"token": token
    	"refresh_token": 用于刷新token
    	"expires_in": token到期时间（以秒为单位）
    	"user_info": 
    	{
    		"uid": 微信号
    		"name": 昵称
		}
	}
}
```

**refresh_token的接口**

```http
/refresh_token
```

**refresh_token 成功后返回**

```json
{
    "code": http状态码
    "message": 信息
    "data":
    {
    	"token": 刷新之后的token
    	"refresh_token": 新的用于刷新的token
    	"expires_in": token到期时间（秒）
	}
}
```



























