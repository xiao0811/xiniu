# 玺牛后台项目

##### 环境:   
`go 1.15(大于 1.13)`   
`docker`   
`mysql`   
`redis`

#### API
<p>所有接口数据请求都是用json传输, POST</p>
e.g:   
 
request:   
```json
{
    "phone": "18949883581",
    "role": 1,
    "password": "123456"
}
```
response:   
```json
{
    "code": 200,
    "message": "ok",
    "data": {
        "id": 6,
        "phone": "18949883581",
        "real_name": "",
        "gender": 0,
        "birthday": "0001-01-01 00:00:00",
        "identification": "",
        "role": 1,
        "marshalling_id": 1,
        "marshalling": {
            "id": 0,
            "name": "",
            "status": 0,
            "type": 0,
            "created_at": "0001-01-01 00:00:00",
            "updated_at": "0001-01-01 00:00:00"
        },
        "created_at": "2020-11-18 20:04:39",
        "updated_at": "2020-11-18 20:04:39"
    }
}
```
