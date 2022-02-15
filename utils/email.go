package utils

import (
	"strconv"

	"github.com/phoenix-next/phoenix-server/global"
	"gopkg.in/gomail.v2"
)

func SendMail(mailTo []string, subject string, body string) error {
	// 定义邮箱服务器连接信息，如果是网易邮箱 password 填密码，qq邮箱填授权码
	// 具体信息可在 phoenix-config.yml 中填写
	//	mailConn := map[string]string{
	//  	"user": "xxx@163.com",
	//  	"password": "your password",
	//  	"host": "smtp.163.com",
	//  	"port": "465",
	//	}

	mailConn := map[string]string{
		"user":     global.VP.GetString("email.user"),
		"password": global.VP.GetString("email.password"),
		"host":     global.VP.GetString("email.host"),
		"port":     global.VP.GetString("email.port"),
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], "Slime scholar")) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["password"])

	err := d.DialAndSend(m)
	return err

}

// SendRegisterEmail 发送注册验证码邮件
func SendRegisterEmail(email string, number int) {
	// 邮件标题
	subject := "欢迎注册PhoeniX学习平台"
	// 邮件接收方
	mailTo := []string{
		email,
	}
	// 邮件内容
	body := `您好，欢迎注册PhoeniX学习平台！
			 本平台旨在提高题目评测速度，提供更优质的知识共享与交流社区，帮助您更轻松地掌握新技能。
			 您的注册验证码为：` + strconv.Itoa(number)
	// 发送邮件
	err := SendMail(mailTo, subject, body)
	if err != nil {
		global.LOG.Panic("SendRegisterEmail: send email error")
	}
}
