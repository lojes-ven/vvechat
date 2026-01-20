## 好友申请

### 加载好友申请列表（http）

```http
GET /api/auth/friendship_requests
Authorization: Bearer <access_token>
```

成功返回：

```json
{
    "code": 200,
    "message": "success",
    "data": [
        {
            "request_id": "123",
            "sender_id": "456",
            "sender_name": "李四",
            "verification_message": "我是你的大学同学",
            "status": "pending",
            "created_at": "2026-01-15T09:30:00Z"
        }
    ]
}
```

失败返回示例：

```json
{
    "code": 500,
    "message": "服务器错误"
}
```

### 发送好友申请（http）

```http
POST /api/auth/friendship_requests
Content-Type: application/json
Authorization: Bearer <access_token>
```

请求体：

```json
{
    "receiver_id": "67890",
    "sender_name": "我的昵称",
    "verification_message": "你好，我们加个好友吧"
}
```

成功返回：

```json
{
    "code": 201,
    "message": "发送成功",
    "data": null
}
```

**websocket:**

- 推送事件类型：

  `new_friend_request`

- 推送消息 JSON 格式：

```json
{
  "type": "new_friend_request",
  "data": {
    "sender_id": 12345,
    "sender_name": "张三",
    "message": "你好，我们加个好友吧",
    "request_id": 67890,
    "created_at": "2026-01-15T09:30:00Z"
  }
}
```

字段说明：
- `type`：事件类型，固定为 `new_friend_request`。

- `data.sender_id`：发送者的用户 ID。

- `data.sender_name`：发送者在申请中填写的名字/昵称。

- `data.message`：验证/附言文本（即 `verification_message`）。

- `data.request_id`：服务端创建的好友申请记录 ID（用于前端在接收通知后拉取或一键跳转到申请详情）。

- `data.created_at`：申请创建时间（UTC 时间戳格式）。

  

- 离线处理：

  如果接收者不在线，WebSocket 推送不会生效。接收者仍然可以通过之前的 HTTP 接口 `GET /api/auth/friendship_requests` 拉取未处理的申请。

### 同意好友申请（http）

```http
POST /api/auth/friendship_requests/{request_id}
Authorization: Bearer <access_token>
```

成功返回：

```json
{
    "code": 201,
    "message": "success",
    "data": null
}
```

失败返回示例：

```json
{
    "code": 400,
    "message": "requestID错误"
}
```

### 删除好友申请（http）

```http
DELETE /api/auth/friendship_requests/{request_id}
Authorization: Bearer <access_token>
```

成功返回：

```json
{
    "code": 201,
    "message": "success",
    "data": null
}
```

---

## 好友关系

### 加载好友列表（http）

```http
GET /api/auth/friendships
Authorization: Bearer <access_token>
```

成功返回：

```json
{
    "code": 200,
    "message": "success",
    "data": [
        {
            "friendship_id": "12345",
            "friend_id": "67890",
            "friend_remark": "同事"
        }
    ]
}
```

### 删除好友（http）

```http
DELETE /api/auth/friendships/{friend_id}
Authorization: Bearer <access_token>
```

成功返回：

```json
{
    "code": 201,
    "message": "success",
    "data": null
}
```

失败返回示例：

```json
{
    "code": 400,
    "message": "好友不存在"
}
```

### 修改好友备注（http）

```http
POST /api/auth/friendships/remark/{friend_id}
Authorization: Bearer <access_token>
Content-Type: application/json
```

请求体：

```json
{
    "remark": "大学同学"
}
```

成功返回：

```json
{
    "code": 200,
    "message": "success",
    "data": null
}
```
