![NekoCAS](https://img.cdn.n3ko.co/lsky/2020/04/10/c33bfa9cfc5b9.png)

中央认证服务 / Central Authentication Service

## 搭建开发环境

NekoCAS 需要以下依赖：
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) (v1.8.3 or higher)
- [Go](https://golang.org/doc/install) (v1.14 or higher)

### 获取代码

```
git clone --depth 1 https://github.com/NekoWheel/NekoCAS
```

### 安装项目依赖

```bash
go mod tidy
```

### 项目说明
NekoCAS 为 NekoWheel 的中央认证服务，注意其**只实现**了 CAS Protocol 中的如下接口：
* `/login` v1
* `/logout` v1
* `/validate` v1
* `/serviceValidate` v2
