package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/larspensjo/config"
	"guard_dog/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	configFile = flag.String("configfile", "config.ini", "General configuration file")

	port         string
	notifyUrl    string
	intervalTime int
	msgKey       string
	at           string

	monitorings = make(map[string]string)

	pushRecord = make(map[string]time.Time)
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	
	fmt.Println("======= 欢迎使用服务器看门狗 =======")
	fmt.Println("======= ======================== =======")

	//显示公司logo 延迟2秒
	time.Sleep(2 * time.Second)
	fmt.Println("..")
	fmt.Println("....")
	fmt.Println("......")

	initConfig()

	go startCheck()

	sendPushMsg("看门狗启动完毕........")

	startServer()

}

func startCheck() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			sendPushMsg("服务器状态通知出现bug，现在已经恢复，请检查日志。")
			//发现异常后 继续调用自己
			startCheck()
		}
	}()

	texts := make([]string, 10)

	for {
		texts = append([]string{})

		var check =  false

		var interval time.Duration
		interval = time.Duration(intervalTime)
		time.Sleep(interval * time.Second)
		for k,v := range monitorings {
			check = utils.TcpStatusCheck(v)
			fmt.Printf(k + " >> check >> check %t \n", check)
			if !check {
				texts = append(texts, k)
			}
		}

		text := "请注意以下服务出现异常  "

		for _,s := range texts {
			text += "," + s
		}

		text += " !!!"
		t, ok := pushRecord[text]
		if ok {
			tt := time.Now().Sub(t)
			if tt.Minutes() < 60 {
				fmt.Println("1小时内重复推送内容 >> " + text)
				continue
			}

		}

		if len(texts) > 0 {
			fmt.Println("=====================")
			ret,_ := sendPushMsg(text)
			if ret {
				pushRecord[text] = time.Now()
			}
		}

		data := time.Now().Format("2006-01-02 15:04:05")

		fmt.Println(data + " >> 监测服务运行正常....")
	}
}

//启动 web service 服务
func startServer() {

	mux := http.NewServeMux()
	//定义简单http 请求方法
	mux.HandleFunc("/send_push_msg", pushMsg)

	server := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 10, //设置3秒的写超时
		Handler:      mux,
	}

	fmt.Println("程序已经启动 消息推送地址 http://127.0.0.1:18080/send_push_msg")
	fmt.Println("curl -X POST http://127.0.0.1:18080/send_push_msg -d 推送消息内容")
	log.Fatal(server.ListenAndServe())
}

func pushMsg(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte("msg push status : failure"))
		return
	}

	if len(body) == 0 {
		w.Write([]byte("msg push status : failure"))
		return
	}

	msg := string(body)

	t, ok := pushRecord[msg]

	ret := true

	if ok {
		tt := time.Now().Sub(t)
		if tt.Minutes() < 60 {
			fmt.Println("1小时内重复推送内容 >> " + msg)
		} else {
			ret, _ := sendPushMsg(msg)
			if ret {
				pushRecord[msg] = time.Now()
			}
		}
	} else {
		ret, _ := sendPushMsg(msg)
		if ret {
			pushRecord[msg] = time.Now()
		}
	}


	if ret {
		w.Write([]byte("msg push status : success"))
	} else {
		w.Write([]byte("msg push status : failure"))
	}
}

func sendPushMsg(msg string) (bool, string) {

	mapmsg := getDefSendMsg(msg)

	bytes, _ := json.Marshal(mapmsg)

	code, retStr := utils.HttpPost(notifyUrl, string(bytes))


	if code == 200 {
		retjson := make(map[string]interface{})

		err := json.Unmarshal([]byte(retStr), &retjson)

		x := retjson["errmsg"]

		v, ok := x.(string)

		if err == nil && ok && v == "ok" {
			return true, retStr
		} else {
			return false, retStr
		}
	} else {
		return false, retStr
	}
}

func getDefSendMsg(textMsg string) map[string]interface{} {

	msg := make(map[string]interface{})
	msg["msgtype"] = "text"

	text := make(map[string]string)
	text["content"] = msgKey + " > " + textMsg

	msg["text"] = text

	if at != "" {
		ats := strings.Split(at, ",")
		mat := make(map[string]interface{})
		mat["atMobiles"] = ats
		mat["isAtAll"] = false
		msg["at"] = mat
	}

	return msg
}

func initConfig() {
	cfg, err := config.ReadDefault(*configFile)

	if err != nil {
		fmt.Println(">>> 关键配置文件查找失败 !!")
		os.Exit(1)
	}

	port, err = cfg.String("app_config", "port")

	if err != nil {
		fmt.Println(">>> app_config port 配置文件读取错误!!")
		os.Exit(1)
	}

	notifyUrl, err = cfg.String("app_config", "notify_url")

	if err != nil {
		fmt.Println(">>> app_config notify_url 配置文件读取错误!!")
		os.Exit(1)
	}

	intervalTime, err = cfg.Int("app_config", "interval_time")

	if err != nil {
		fmt.Println(">>> app_config interval_time 配置文件读取错误!!")
		os.Exit(1)
	}

	msgKey, err = cfg.String("app_config", "msg_key")

	if err != nil {
		fmt.Println(">>> app_config msg_key 配置文件读取错误!!")
		os.Exit(1)
	}

	at, err = cfg.String("app_config", "at")

	if err != nil {
		fmt.Println(">>> app_config at配置文件读取错误!!")
		os.Exit(1)
	}

	options, err := cfg.SectionOptions("monitorings")

	if err != nil {
		fmt.Println(">>> monitorings 配置文件读取错误!!")
		os.Exit(1)
	}

	for _, v := range options {
		vv, err := cfg.String("monitorings", v)
		if err != nil {
			fmt.Println(">>> monitorings 配置文件读取错误!!")
			os.Exit(1)
		}
		monitorings[v] = vv
	}

	fmt.Println("..........系统初始化 完毕.")

	for k,v := range monitorings {
		fmt.Println("监控的服务信息 >>> " + k + " " + v)
	}

}
