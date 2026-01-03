- **打开微信时， 加载会话列表**

  加载会话列表前端接口：

  ```http
  GET api/auth/conversations
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  ```

  失败后端返回：

  ```json
  
  ```

- **会话详情（表现为用户点击进入聊天窗口）**

  前端接口

  ```http
  GET api/auth/conversations/{conversation_id}
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  {
      "code":  200,
      "message": 信息
  }
  ```

  失败后端返回：

  ```json
  {
      "code":  400 / 401 / 500 / 409,
      "message": 信息
  }
  ```

- **发送消息**

  前端接口

  ```http
  POST /api/auth/messages
  Content-Type: application/json

  {
  	"receiver_id": "接收者ID",
  	"content": "消息内容"
  }
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  {
      "code":  200,
      "message": "success",
      "data": "消息ID"
  }
  ```

  失败后端返回：

  ```json
  {
      "code":  400,
      "message": "错误信息"
  }
  ```

- **撤回消息**

  前端接口

  ```http
  DELETE /api/auth/messages/recall
  Content-Type: application/json

  {
  	"message_id": "消息ID"
  }
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  {
      "code":  200,
      "message": "success"
  }
  ```

  失败后端返回：

  ```json
  {
      "code":  400,
      "message": "错误信息"
  }
  ```

- **删除消息**

  前端接口

  ```http
  DELETE /api/auth/messages/delete
  Content-Type: application/json

  {
  	"message_id": "消息ID"
  }
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  {
      "code":  200,
      "message": "success"
  }
  ```

  失败后端返回：

  ```json
  {
      "code":  400,
      "message": "错误信息"
  }
  ```

  





















