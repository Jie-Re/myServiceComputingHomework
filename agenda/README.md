# 安装cobra并完成小案例
## 什么是cobra
cobra： 一个用生成程序应用和命令行文件的程序
具体可参见[golang命令行库cobra的使用 ](https://www.cnblogs.com/borey/p/5715641.html)

## 项目目录创建
可参考[我之前的博客](https://blog.csdn.net/xxiangyusb/article/details/100858000)及[潘老师的课程网站](https://pmlpml.github.io/ServiceComputingOnCloud/ex-install-go)
## 安装cobra
同样由于`https://golang.org`的访问问题，我们首先需要下载其中cobra需要的依赖
- 下载依赖
	
	```bash
	$ cd $GOPATH/src/github.com/golang
	$ git clone https://github.com/golang/sys.git
	$ git clone https://github.com/golang/text.git
	```
- 克隆完成后，复制到相应目录

	```bash
	$ cp $GOPATH/src/github.com/golang/sys $GOPATH/src/golang.org/x/ -rf
	$ cp $GOPATH/src/github.com/golang/text $GOPATH/src/golang.org/x/ -rf
	```

- 编译cobra

	```go
	$ go get -v github.com/spf13/cobra/cobra
	```

- 安装cobra

	```go
	$ go install github.com/spf13/cobra/cobra
	```
- 安装后，在`$GOPATH/bin`下出现了可执行程序
![cobra-installed](https://img-blog.csdnimg.cn/20191029150409829.PNG#pic_center)
## 使用cobra生成应用程序
例如，
- 创建一个基于CLI的命令程序agenda：

	```bash
	$ cd $GOPATH/src
	$ cobra init agenda --pkg-name agenda
	```
![cobra-init-directory](https://img-blog.csdnimg.cn/20191029152929419.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3h4aWFuZ3l1c2I=,size_16,color_FFFFFF,t_70#pic_center)
- 添加`register.go`文件

	```bash
	$ cd agenda
	$ cobra add register
	```
![tree-agenda](https://img-blog.csdnimg.cn/20191029153144940.PNG#pic_center)
- 修改`register.go`的代码
	- `init()`函数中添加如下内容：

		```go
		registerCmd.Flags().StringP("user", "u", "Anonymous", "Help message for username")
		```

	- `Run`匿名回调函数中添加：

		```go
				username, _ := cmd.Flags().GetString("user")
				fmt.Println("register called by " + username)
		```
- 测试：

	```bash
	$ go run main.go register --user=wyb
	```
![register-test](https://img-blog.csdnimg.cn/20191029153957417.PNG#pic_center)
![register-test2](https://img-blog.csdnimg.cn/2019102915442740.PNG#pic_center)
# agenda开发实战
通过上面测试案例，我们已能感受到利用cobra开发命令行程序的方便性。由于本次作业只要求实现两条命令，这也就大大降低了开发难度。
## Agenda业务需求
Agenda主要有以下两个功能
- 用户注册
	- 注册新用户时，用户需设置一个唯一的用户名和一个密码。另外，还需登记邮箱及电话信息
	- 如果注册时提供的用户名已由其他用户使用，应反馈一个适当的出错信息；成功注册后，亦应反馈一个成功注册的信息
- 用户登录
	- 用户使用用户名和密码登录 Agenda 系统
	- 用户名和密码同时正确则登录成功并反馈一个成功登录的信息。否则，登录失败并反馈一个失败登录的信息
## agenda程序设计
`code root.go`，修改变量`rootCmd`中的`Short`和`Long`的描述性内容。
另外，viper是cobra集成的配置文件读取的库，这里不需要使用，我们可以注释掉（不注释可能生成的应用程序很大）
## 用户注册
### 添加子程序
```bash
$ cd agenda
$ cobra add register
```
### 设置命令行参数
`code register.go`，修改变量`registerCmd`中的`Short`和`Long`的描述性内容，在`init()`函数中设置命令行参数

```go
	registerCmd.Flags().StringP("user", "u", "", "Help message for username")
	registerCmd.Flags().StringP("password", "p", "", "Help message for password")
	registerCmd.Flags().StringP("email", "e", "", "Help message for email")
```
### 检查参数合理性
```go
	username, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	emailaddr, _ := cmd.Flags().GetString("email")
	if len(username) == 0 {
		fmt.Println("Error: Username must be set")
		cmd.Help()
		return
	}
	if len(password) == 0 {
		fmt.Println("Error: Password must be set")
		cmd.Help()
		return
	}
	if len(emailaddr) == 0 {
		fmt.Println("Error: Email address must be set")
		cmd.Help()
		return
	}
	//检查邮箱地址合理性
	matched, _ := regexp.MatchString(`[\w-]+@[\w]+(?:\.[\w]+)+`, emailaddr)
	if matched == false {
		fmt.Println("Error: Your email address is invalid, please check")
		return
	}
```
**友情提示**：检查邮箱地址合理性时可直接借助`regexp`包（需import），具体实现：
```go
matched, _ := regexp.MatchString(`[\w-]+@[\w]+(?:\.[\w]+)+`, emailaddr)
```
### 检查用户名唯一性
```go
//check if username was unique
if fileObjR, errR := os.OpenFile("users.txt", os.O_RDONLY|os.O_CREATE, 0644); errR == nil {
	defer fileObjR.Close()
	if contents, err := ioutil.ReadAll(fileObjR); err == nil {
		result := strings.Replace(string(contents), "\n", "", 0)
		infos := strings.Split(result, "\n")
		for i := 0; i < len(infos); i += 3 {
			if username == infos[i] {
				fmt.Println("Error: This username has been used, please choose another one")
				return
			}
		}
	}
}
```

## 用户登录
### 添加子程序
```bash
$ cd agenda
$ cobra add login
```
### 设置命令行参数
`code login.go`，修改变量`loginCmd`中的`Short`和`Long`的描述性内容，在`init()`函数中设置命令行参数

```go
	loginCmd.Flags().StringP("user", "u", "", "Help message for username")
	loginCmd.Flags().StringP("password", "p", "", "Help message for password")
```
### 检查参数合理性

```go
	username, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	if len(username) == 0 {
		fmt.Println("Error: Username must be set")
		cmd.Help()
		return
	}
	if len(password) == 0 {
		fmt.Println("Error: Password must be set")
		cmd.Help()
		return
	}
```
### 检查用户名是否存在及密码正确性

```go
//check if password for username was correct
if fileObjR, errR := os.OpenFile("users.txt", os.O_RDONLY|os.O_CREATE, 0644); errR == nil {
	defer fileObjR.Close()
	if contents, err := ioutil.ReadAll(fileObjR); err == nil {
		result := strings.Replace(string(contents), "\n", "", 0)
		infos := strings.Split(result, "\n")
		for i := 0; i < len(infos); i += 3 {
			if username == infos[i] {
				if password == infos[i+1] {
					fmt.Println("Succeed:  user " + username + " login successfully")
					return
				}
				fmt.Println("Fail: password incorrect, please try again")
				return
			}
		}
		fmt.Println("Fail: you are not registered, please register first")
	}
}
```
## 安装agenda程序

```bash
$ go install agenda
```
## 运行测试
1. `agenda help`
![agenda-help](https://img-blog.csdnimg.cn/20191029182638566.PNG?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3h4aWFuZ3l1c2I=,size_16,color_FFFFFF,t_70#pic_center)
2. `agenda register -h`
![register-help](https://img-blog.csdnimg.cn/20191029182822278.PNG?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3h4aWFuZ3l1c2I=,size_16,color_FFFFFF,t_70#pic_center)
3. `agenda login -h`
![login-help](https://img-blog.csdnimg.cn/20191029183105363.PNG?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3h4aWFuZ3l1c2I=,size_16,color_FFFFFF,t_70#pic_center)
4. `agenda register -u wyb -p 0805 -e 19970805@163.com`
![register-success](https://img-blog.csdnimg.cn/20191029183144620.PNG#pic_center)  
**相应错误测试**：  
未设置用户名`agenda register -p 0805 -e 19970805@163.com`  
![register-error1](https://img-blog.csdnimg.cn/20191029183402929.PNG?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3h4aWFuZ3l1c2I=,size_16,color_FFFFFF,t_70#pic_center)  
	未设置密码`agenda register -u web -e 19970805@163.com`  
![register-error2](https://img-blog.csdnimg.cn/20191029183450831.PNG?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3h4aWFuZ3l1c2I=,size_16,color_FFFFFF,t_70#pic_center)  
	未设置邮箱`agenda register -u web -p 0805`  
![register-error3](https://img-blog.csdnimg.cn/20191029183544755.PNG?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3h4aWFuZ3l1c2I=,size_16,color_FFFFFF,t_70#pic_center)  
	邮箱地址不合法`agenda register -u web -p 0805 -e 19970805`  
![register-error4](https://img-blog.csdnimg.cn/20191029183630771.PNG#pic_center)  
	用户名重复`agenda register -u wyb -p 1997 -e 0805@163.com`  
![register-error5](https://img-blog.csdnimg.cn/20191029183714800.PNG#pic_center)  
5. `agenda login -u wyb -p 0805`
![login](https://img-blog.csdnimg.cn/20191029184031307.PNG)  
	**相应错误测试**：  
	用户名不存在`agenda login -u aaa -p 0805`  
![login-fail1](https://img-blog.csdnimg.cn/20191029184118670.PNG)  
	密码错误`agenda login -u wyb -p 1997`  
	![login-fail2](https://img-blog.csdnimg.cn/2019102918423852.PNG)  
# 本文博客链接
[服务计算Agend小程序开发作业心得](https://blog.csdn.net/xxiangyusb/article/details/102799076)
