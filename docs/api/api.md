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
      "table_id": 1,
      "table_name": "标签名称"
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
  "category_id": 1,
  "private_tables": [
    {
      "table_id": 1,
      "columns": [
        {
          "id": 1,
          "value": ["1","2"] # 多个枚举值则传入多个字符串
        },
        {
          "id": 2,
          "value": ["(1,2)"] # 多个范围则多传多个范围字符串
        }
      ]
    }
  ],
  "public_tables": [
    {
      "table_id": 2,
      "columns": [
        {
          "id": 4,
          "value": ["1","2"]
        }
      ]
    }
  ],
  "portrait_tables": [
    {
      "table_id": 2
    }
  ]
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


```json
{
  "parent_node_tag_list": [
    {
      "table_name": "dmp_tagfactory_detail_fmcg_crowd",
      "filter": "label!='0'",
      "field": "label,label2,label3"
    },
    {
      "table_name": "dmp_tagfactory_detail_alipay_amt_level_n",
      "filter": "alipay_amt_level>=0",
      "field": "alipay_amt_level"
    }
  ],
  "child_columns_set": [
    {
      "tag_name": "t1_label",
      "tag_value": [
        "资深白领",
        "小镇青年",
        "新锐白领",
        "小镇中年",
        "都市蓝领",
        "GenZ",
        "精致妈妈"
      ],
      "tag_type": "enumerate"
    },
    {
      "tag_name": "t1_label2",
      "tag_value": [
        "资深白领",
        "小镇青年",
        "新锐白领",
        "小镇中年",
        "都市蓝领",
        "GenZ",
        "精致妈妈"
      ],
      "tag_type": "enumerate"
    },
    {
      "tag_name": "t2_alipay_amt_level",
      "tag_value": [
        [
          1,
          2
        ],
        [
          3,
          5
        ]
      ],
      "tag_type": "range"
    }
  ],
  "portrait_table_list": [
    {
      "table_name": "dmp_tagfactory_detail_gprofile_gender_n",
      "filter": "gender>=0",
      "field": "gender"
    },
    {
      "table_name": "dmp_tagfactory_detail_gprofile_age_n",
      "filter": "age>0",
      "field": "age"
    },
    {
      "table_name": "dmp_tagfactory_detail_alipay_amt_level_n",
      "filter": "alipay_amt_level>0",
      "field": "alipay_amt_level"
    }
  ],
  "parent_private_table_name": "dmp_tagfactory_similar_shop_behavior_s4_n",
  "private_columns_set": [
    {
      "tag_name": "type",
      "tag_value": [
        [
          1,
          2
        ],
        [
          3,
          4
        ]
      ],
      "tag_type": "range"
    },
    {
      "tag_name": "level",
      "tag_value": [
        1,
        2
      ]
    },
    {
      "tag_name": "period",
      "tag_value": [
        30
      ]
    }
  ]
}
```


```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "uid": "123456",
    "name": "",
    "age": "",
    "institution": "",
    "topic": "",
    "gender": 1
  }
}
```