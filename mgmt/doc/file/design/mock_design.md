##### 支持的Mock路径由
- http://{IP}:{PORT}/mock/file/name.xml 可返回上传文件下name.xml的内容，原样返回
- http://{IP}:{PORT}/mock/data/slow?sleep=N  N为整数，表示等待的时间，模拟真实三方请求返回有延迟的情况，可返回JSON格式的数据
- http://{IP}:{PORT}/mock/data/quick  可返回JSON格式的数据，内容如下

##### 支持的样例信息，各项信息是随机的：
```
{
    addr: "XX市XX市XX路XX号XX小区XX单元XX室",
    age: "XX",
    bankCard: "XXXXX",
    bool: "XX",
    grade: "X",
    id_no: "XXXXX",
    mobile: "XX",
    name: "XXXXX",
    sex: "XXX",
    weight: "XX",
    yesOrNo: "XX"
}
```
