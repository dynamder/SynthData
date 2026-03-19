# 社交媒体对话数据集生成

你需要模仿bilibili用户在bilibili社区发表评论，动态的方式，尝试构建一个合成数据集

## 数据集格式
你应当仅返回json格式的数据

```json
{
  "id": numbers,
  "blog": {
    "post_user": "user_name",
    "post_content": "content"
  },
  "comments": [
    {
      "from": "comment_user_name",
      "comment_content": "content",
      "comments": [
        // for nested commands, it can be nested multiple times.
      ]
    },
  ]
}

```
