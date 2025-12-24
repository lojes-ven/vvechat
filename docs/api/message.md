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

- **私发消息**

  前端接口

  ```http
  POST api/auth/messages
  json
  {
  	"conversation_id": 给哪个会话发
  	"content": 消息内容
  }
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  {
      "code":  200,
      "message": 信息,
      "data": message_id（信息表的主键）
  }
  ```

  失败后端返回：

  ```json
  {
      "code":  400 / 401 / 500 / 409,
      "message": 信息
  }
  ```

- **撤回消息**

  前端接口

  ```http
  DELETE api/auth/message/recall
  {
  	"message_id":消息表的主键
  }
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  {
      "code":  200,
      "message": 信息,
  }
  ```

  失败后端返回：

  ```json
  {
      "code":  400 / 401 / 500 / 409,
      "message": 信息
  }
  ```

- **删除消息**

  ```http
  DELETE api/auth/message/delete
  {
  	"message_id":消息表的主键
  }
  Authorization: Bearer <token>
  ```

  成功后端返回：

  ```json
  {
      "code":  200,
      "message": 信息,
  }
  ```

  失败后端返回：

  ```json
  {
      "code":  400 / 401 / 500 / 409,
      "message": 信息
  }
  ```

  





















