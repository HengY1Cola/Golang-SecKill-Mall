package main

import (
	"Iris/common"
	"Iris/datamodels"
	"Iris/rabbitmq"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var localHost = ""
var hostArray = []string{"127.0.0.1", "127.0.0.1"} //设置集群地址，最好内外IP
var port = "8081"
var hashConsistent *common.Consistent
var rabbitMqValidate *rabbitmq.RabbitMQ
var GetOneIp = "127.0.0.1" //数量控制接口服务器内网IP，或者getOne的SLB内网IP
var GetOnePort = "8084"
var interval = 20

// BlackList 黑名单
type BlackList struct {
	listArray map[int]bool
	sync.RWMutex
}

// AccessControl 用来存放控制信息，
type AccessControl struct { //用来存放用户想要存放的信息
	sourcesArray map[int]time.Time
	sync.RWMutex
}

//创建全局变量
var accessControl = &AccessControl{sourcesArray: make(map[int]time.Time)}
var blackList = &BlackList{listArray: make(map[int]bool)}

// GetNewRecord 获取制定的数据
func (m *AccessControl) GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.sourcesArray[uid]
}

// SetNewRecord 设置记录
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourcesArray[uid] = time.Now()
}

// GetDistributedRight 去对应的主机去判断是否通过
func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	//todo 获取用户UID
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	// todo 根据用户ID，判断获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	//todo 判断是否为本机
	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value) //执行本机数据读取和校验
	} else {
		return GetDataFromOtherMap(hostRequest, req) //不是本机充当代理访问数据返回结果
	}
}

// GetDataFromMap 获取本机map，并且处理业务逻辑，返回的结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	// todo 拿到uid
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	// todo 判断是否在黑名单
	if blackList.GetBlackList(uidInt) {
		return false
	}
	// todo 判断是否在有效时间内 interval决定多少秒还可以买
	data := m.GetNewRecord(uidInt)
	if !data.IsZero() {
		if data.Add(time.Duration(interval) * time.Second).After(time.Now()) { // 距离上次还没到20s
			return false
		}
	}
	m.SetNewRecord(uidInt)
	return true
}

// GetBlackList 获取黑名单
func (b *BlackList) GetBlackList(uid int) bool {
	b.Lock()
	defer b.Unlock()
	return b.listArray[uid]
}

// SetBlackList 获取黑名单
func (b *BlackList) SetBlackList(uid int) {
	b.Lock()
	defer b.Unlock()
	b.listArray[uid] = true
}

// GetCurl 模拟请求
func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	// todo 获取Uid/sign
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}
	// todo 模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	// todo 获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

// GetDataFromOtherMap 获取其它节点处理结果
func GetDataFromOtherMap(host string, request *http.Request) bool {
	// todo 请求
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, request)
	if err != nil {
		return false
	}
	// todo 判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// Check 执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	// todo 拿到productID
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	// todo 获取用户cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	// todo 分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}
	// todo 获取数量控制权限，防止秒杀出现超卖现象
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	// todo 判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			// todo 整合下单
			productID, err := strconv.ParseInt(productString, 10, 64) // 获取商品ID
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64) // 获取用户ID
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			message := datamodels.NewMessage(userID, productID) //创建消息体
			byteMessage, err := json.Marshal(message)           //类型转化
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			err = rabbitMqValidate.PublishSimple(string(byteMessage)) //生产消息
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
}

// Auth 统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	err := CheckUserInfo(r) //添加基于cookie的权限验证
	if err != nil {
		return err
	}
	return nil
}

// CheckUserInfo 身份校验函数
func CheckUserInfo(r *http.Request) error {
	// todo 获取Uid，cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户UID Cookie 获取失败！")
	}
	// todo 获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie 获取失败！")
	}
	// todo 对信息进行解密
	signByte, err := common.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改！")
	}
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败！")
}

//自定义逻辑判断
func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	// todo 负载均衡器设置
	hashConsistent = common.NewConsistent()
	for _, v := range hostArray { //采用一致性hash算法，添加节点
		hashConsistent.Add(v)
	}
	// todo 获取本机IP
	localIp, err := common.GetEntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println("本机IP:", localHost)

	// todo 消息队列
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("Product")
	defer rabbitMqValidate.Destroy()

	// todo 过滤器
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth) //注册拦截器
	filter.RegisterFilterUri("/checkRight", Auth)
	// todo 启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))
	http.ListenAndServe(":8083", nil)
}
