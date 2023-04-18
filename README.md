# MiDa Common Tool

觅搭公共组件

发布新版本方法

```bash
git tag v0.1.0
git push origin v0.1.0
```

在其它项目中更新或者使用 mida-tool

先设置环境变量

```bash
export GOPRIVATE=github.com/letjoy-club/mida-tool
```

编辑 ~/.gitconfig 文件，让 golang 库的请求优先走 ssh 认证。

```
[url "git@github.com:"]
	insteadOf = https://github.com/
```


加载

```bash
go get github.com/letjoy-club/mida-tool
```
