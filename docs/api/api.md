# 对接爬虫接口
## /api/v1/dmp/auth/complete (POST)
### desc
商户认证完成
### request
```json
{
  "merchant_name": "商户名称",
  "merchant_id": "商户ID"
}
```

### response
```json
{
  "code": 0,
  "msg": "success"
}
```

## /api/v1/dmp/categorys (GET)
### desc
获取分类列表
### response
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "cate_id": "分类ID",
      "cate_name": "分类名称",
      "cate_level": "分类级别"
    }
  ]
}
```

# 客户前台接口
## /api/v1/dmp/cal/result (GET)
### desc
获取计算结果
### request
```json
{}
```

### response
```json
{
    "code": 0,
    "msg": "success",
    "data": {
      "cal_results": [
        {
          "cal_time": 123456789,
          "cate_id": "分类ID",
          "cate_name": "分类名称"
        }
      ]
    }
}
```

## /api/v1/dmp/store (GET)
### desc
获取店铺信息，获取该用户下绑定的店铺
### request
```json
{}
```
### response
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "store_id": "店铺ID",
      "store_name": "店铺名称",
      "store_type": "店铺类型"
    }
  ]
}
```

## /api/v1/dmp/store/add (POST)
### desc
添加店铺
### request
```json
{
  "store_name": "店铺名称",
  "store_type": "店铺类型"
}
```
### response
```json
{
  "code": 0,
  "msg": "success"
}
```

## /api/v1/dmp/vip (GET)
### desc
获取套餐信息
### request
```json
{}
```
### response
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "vip_id": "套餐ID",
      "vip_name": "套餐名称",
      "vip_price": "套餐价格",
      "vip_desc": "套餐描述",
      "vip_yearly_price": "套餐年度价格"
    }
  ]
}
```

# 管理后台接口
## /api/v1/dmp/admin/tag/list (GET)
### desc
获取标签配置
### request
```json
{}
```
### response
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "tag_id": "标签ID",
      "tag_name": "标签名称",
      "tag_desc": "标签描述",
      "tag_columns": [
        {
          "column_id": "字段ID",
          "column_name": "字段名称",
          "column_type": "字段类型"
        }
      ]
    }
  ]
}
```

## /api/v1/dmp/admin/tag/config (POST)
### desc
配置标签
### request
```json
{
  "cate_id": "分类ID",
  "private_tags": {
    "tag_id": "标签ID",
    "tag_columns": [
      {
        "column_id": "字段ID",
        "column_value": "字段取值",
        "column_type": "字段类型"
      }
    ]
  },
  "public_tags": {
    "tag_id": "标签ID",
    "tag_columns": [
      {
        "column_id": "字段ID",
        "column_value": "字段取值",
        "column_type": "字段类型"
      }
    ]
  },
  "user_tags": {
    "tag_id": "标签ID",
    "tag_columns": [
      {
        "column_id": "字段ID",
        "column_value": "字段取值",
        "column_type": "字段类型"
      }
    ]
  }
}
```

## /api/v1/dmp/admin/relation/list (GET)
### desc
获取关系配置
### request
```json
{}
```

### response
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "cate_id": "分类ID",
      "private_tags": {
        "tag_id": "标签ID",
        "tag_columns": [
          {
            "column_id": "字段ID",
            "column_value": "字段取值",
            "column_type": "字段类型"
          }
        ]
      },
      "public_tags": {
        "tag_id": "标签ID",
        "tag_columns": [
          {
            "column_id": "字段ID",
            "column_value": "字段取值",
            "column_type": "字段类型"
          }
        ]
      },
      "user_tags": {
        "tag_id": "标签ID",
        "tag_columns": [
          {
            "column_id": "字段ID",
            "column_value": "字段取值",
            "column_type": "字段类型"
          }
        ]
      }
    }
  ]
}
```

## /api/v1/dmp/admin/vip (POST)
### desc
商户套餐管理
### request
```json
{
  "user_id": "用户ID",
  "vip_id": "套餐ID",
  "end_time": 11111
}
```

### response
```json
{
  "code": 0,
  "msg": "success"
}
```

## /api/v1/dmp/admin/user/store/bind (POST)
### desc
用户和商铺绑定
### request
```json
{
  "user_id": "用户ID",
  "store_id": "商铺ID"
}
```
### response
```json
{
  "code": 0,
  "msg": "success"
}
```