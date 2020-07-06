package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"gopkg.in/go-playground/webhooks.v5/github"
	_ "gopkg.in/go-playground/webhooks.v5/github"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http"
	"time"
)

const (
	path = "/"
)

func main() {
	hook, _ := github.New(github.Options.Secret("ImJw"))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			rebuild()
		}
	})
	http.ListenAndServe(":80", nil)

}



func rebuild(){
	sshHost := "_"
	sshUser := "_"
	sshPassword := "_"
	sshType := "_"//password 或者 key
	sshKeyPath := ""//ssh id_rsa.id 路径"
	sshPort := 27947


	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second,//ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	}



	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("创建ssh client 失败",err)
	}
	defer sshClient.Close()


	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败",err)
	}
	defer session.Close()
	//执行远程命令
	combo,err := session.CombinedOutput("cd privateServer; bash ./rebuild.sh &")
	if err != nil {
		log.Fatal("远程执行cmd 失败",err)
	}
	log.Println("命令输出:",string(combo))

}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}