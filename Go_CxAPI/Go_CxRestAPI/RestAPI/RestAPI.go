package RestAPI

import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type CxRestAPI struct {
	ServerIP     string
	URL          string
	ProjectID    string
	ScanID       string
	CxToken      string
	CxUser       string
	CxPassword   string
	CxMethod     string
	ZipFilePath  string
	CxJson       map[string]interface{}
	CxConfig     map[string]interface{}
	status       interface{}
	Body         string
	ScanStatus   string
	RegisterName string
}

type config struct {
	IP       string
	User     string
	Password string
}

//读取url.json文件
func (Readjson *CxRestAPI) Readjson(filepath string) (file map[string]interface{}) {
	Cxjson, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}
	//获取数据
	var jsonfile map[string]interface{}
	er := json.Unmarshal(Cxjson, &jsonfile)
	if er != nil {
		log.Fatal(er)
	}
	//Readjson.CxJson = url_method
	return jsonfile
}

//写入文件
func WriteFile(data []byte) {
	fp, err := os.OpenFile("Go_CxRestAPI/etc/config.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

//获取url和method
func (geturl_met *CxRestAPI) Geturl_method(met string) {
	url_method := geturl_met.CxJson
	str := url_method[met]
	sum := str.(map[string]interface{})
	geturl_met.URL = sum["url_suffix"].(string)
	geturl_met.CxMethod = sum["http_method"].(string)
}

//发送请求
func (SendRequest *CxRestAPI) SendRequest(parameters string) (status int, Body string) {
	//URL拼接
	URL := "http://" + SendRequest.ServerIP + "/cxrestapi" + SendRequest.URL + SendRequest.ProjectID
	Cliet := &http.Client{}
	request, err := http.NewRequest(SendRequest.CxMethod, URL, strings.NewReader(parameters))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", SendRequest.CxToken)
	resp, err := Cliet.Do(request)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//判断是否是获取扫描报告
	if SendRequest.URL == "/reports/sastScan" {
		//报告写入
		er := ioutil.WriteFile(SendRequest.RegisterName, body, 0644)
		if er != nil {
			panic(er)
		}
	}

	//还原ProjectID
	SendRequest.ProjectID = ""
	return resp.StatusCode, string(body)
}

//登陆
func (login *CxRestAPI) Login(JsonFile string) (statu int, resp string) {
	login.Geturl_method(JsonFile)
	parameters := "userName=" + login.CxUser + "&" + "password=" + login.CxPassword
	//调用类里的方法
	status, response := login.SendRequest(parameters)
	return status, response

}

//获取认证token
func (GetToken *CxRestAPI) GetCxToken(JsonFile string) (statu int) {
	GetToken.Geturl_method(JsonFile)
	parameters := "userName=" + GetToken.CxUser + "&" + "password=" + GetToken.CxPassword + "&" + "grant_type=password&" +
		"scope=sast_rest_api&client_id=resource_owner_client&client_secret=014DF517-39D1-4453-B7B3-9930C563627C"
	status, body := GetToken.SendRequest(parameters)
	var token string
	var Body map[string]interface{}
	err := json.Unmarshal([]byte(body), &Body)
	if err != nil {
		log.Panic(err)
	}
	if status == 200 {
		token = Body["access_token"].(string)
		GetToken.CxToken += "Bearer " + token
	} else {
		fmt.Println("构建访问token失败")
	}
	return status

}

//获取所有项目详情&按照项目ID查看
func (Getprojects *CxRestAPI) Getprojects(Project_ID, JsonFile string) {
	Getprojects.Geturl_method(JsonFile)
	parameters := ""
	Getprojects.ProjectID = "/" + Project_ID
	status, body := Getprojects.SendRequest(parameters)
	fmt.Println(status, body)
}

//获取teams的信息
func (getteams *CxRestAPI) GetTeams(JsonFile string) {
	getteams.Geturl_method(JsonFile)
	parameters := ""
	status, body := getteams.SendRequest(parameters)
	fmt.Println(status, body)
}

//使用默认配置创建项目
func (CreateProject *CxRestAPI) CreateProject(ProjectName, JsonFile string) (sta int, inf string) {
	//从url.json文件中获取URL及method
	CreateProject.Geturl_method(JsonFile)
	//拼接parameters
	parameters := "name=" + ProjectName + "&owningTeam=00000000-1111-1111-b111-989c9070eb11&isPublic=true"
	status, body := CreateProject.SendRequest(parameters)

	//转换格式string——>map
	var info string
	var Body map[string]interface{}
	err := json.Unmarshal([]byte(body), &Body)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(Body)
	if status == 201 {
		ID := Body["id"].(float64)
		//类型转换float64——>string
		info = strconv.FormatFloat(ID, 'f', -1, 64)
		fmt.Println("创建项目成功，ID为：", info)
	} else if status == 403 {
		info = Body["messageDetails"].(string)
	} else {
		fmt.Println("其他错误，请重新创新项目")
		info = "0"
	}
	return status, info
}

//上传zip包
func (uploadzipfile *CxRestAPI) upload_source_code_zip_file(ProjectID, JsonFile string) (sta int) {
	uploadzipfile.Geturl_method(JsonFile)
	//字符串拼接fmt.Sprintf(qq%s,"123") #qq123
	URL := "http://" + uploadzipfile.ServerIP + "/cxrestapi" + fmt.Sprintf(uploadzipfile.URL, ProjectID)
	//创建表单文件
	//CreateFormFile 用来创建表单，第一个参数是字段名，第二个参数是文件名
	client := &http.Client{}
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("zippedSource", "BookStoreJava.zip")
	if err != nil {
		panic(err)
	}
	// 从文件读取数据，写入表单
	srcFile, err := os.Open("BookStoreJava.zip")
	if err != nil {
		panic(err)
	}
	//srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
	}
	// 发送表单
	contentType := writer.FormDataContentType()
	// 发送之前必须调用Close()以写入结尾行
	writer.Close()
	request, err := http.NewRequest(uploadzipfile.CxMethod, URL, buf)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", uploadzipfile.CxToken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	return response.StatusCode
}

//开始扫描
func (newscan *CxRestAPI) create_new_scan(ProjectID, JsonFile string) (sta int, ScanID string) {
	newscan.Geturl_method(JsonFile)
	parameters := "projectId=" + ProjectID + "&isIncremental=false&isPublic=false&forceScan=True"
	status, body := newscan.SendRequest(parameters)
	var info string
	if status == 201 {
		var Body map[string]interface{}
		err := json.Unmarshal([]byte(body), &Body)
		if err != nil {
			log.Panic(err)
		}
		id := Body["id"].(float64)
		info = strconv.FormatFloat(id, 'f', -1, 64)
	} else {
		fmt.Println("创建扫描失败")
		info = "0"
	}
	return status, info
}

//获取扫描状态
func (scanQueue *CxRestAPI) get_all_scan_details_in_queue(ScanID, JsonFile string) (sta int, ScanStayus, StatusID string) {
	scanQueue.Geturl_method(JsonFile)
	scanQueue.ProjectID = "/" + ScanID
	parameters := ""
	status, body := scanQueue.SendRequest(parameters)
	var info, ScanStatusID string

	if status == 200 {
		var Body map[string]interface{}
		err := json.Unmarshal([]byte(body), &Body)
		if err != nil {
			log.Panic(err)
		}
		info = Body["stage"].(map[string]interface{})["value"].(string)
		ID := Body["stage"].(map[string]interface{})["id"].(float64)
		ScanStatusID = strconv.FormatFloat(ID, 'f', -1, 64)
		//fmt.Printf("队列状态：%s，状态id：%s", info, ScanStatusID)
	} else {
		fmt.Println("获取扫描状态失败")
		info = "0"
		return
	}

	return status, info, ScanStatusID
}

//查看所有扫描项目
func (scanProjects *CxRestAPI) get_sast_scan_details_by_scan_id(ScanID, JsonFile string) {
	scanProjects.Geturl_method(JsonFile)
	scanProjects.ProjectID = "/" + ScanID
	parameters := ""
	status, body := scanProjects.SendRequest(parameters)
	fmt.Println(status, body)
}

//获取报告，注册扫描报告
func (RegisterScanReport *CxRestAPI) register_scan_report(ScanID, FileType, JsonFile string) (sta int, reportId string) {
	RegisterScanReport.Geturl_method(JsonFile)
	parameters := "reportType=" + FileType + "&scanId=" + ScanID
	status, body := RegisterScanReport.SendRequest(parameters)
	fmt.Println(status, body)
	var info string
	if status == 202 {
		var Body map[string]interface{}
		err := json.Unmarshal([]byte(body), &Body)
		if err != nil {
			log.Panic(err)
		}
		ID := Body["reportId"].(float64)
		info = strconv.FormatFloat(ID, 'f', -1, 64)
	} else {
		fmt.Println("注册扫描报告失败")
		info = "0"
		return
	}
	//fmt.Println(info)
	return status, info
}

//获取报告
func (GetReports *CxRestAPI) get_reports_by_id(ReportId, JsonFile string) {
	GetReports.Geturl_method(JsonFile)
	GetReports.ProjectID = "/" + ReportId
	parameters := ""
	status, _ := GetReports.SendRequest(parameters)
	if status == 200 {
		fmt.Println("报告已生成")
	} else {
		fmt.Println("报告获取失败")
	}
}

func RestAPI() {
	var ProjectID, ScanID, ZipFilePath, Body, ScanStatus, CxIP, CxUser, CxPassword, CxUrl, num, RegisterName, login string
	var CxJson, CxConfig map[string]interface{}
	CxToken := ""
	status := "200"
	CxMethod := "GST"
	//初始化对象
	CxAPI := CxRestAPI{CxIP, CxUrl, ProjectID, ScanID, CxToken,
		CxUser, CxPassword, CxMethod, ZipFilePath, CxJson,
		CxConfig, status, Body, ScanStatus, RegisterName}
	//读取url.json文件
	CxAPI.CxJson = CxAPI.Readjson("Go_CxRestAPI/etc/urls.json")
	CxAPI.CxConfig = CxAPI.Readjson("Go_CxRestAPI/etc/config.json")

	fmt.Println("如上次配置登陆请回车,输入其他则重新输入配置：") //192.168.40.230
	fmt.Scanln(&login)
	if login == "" {
		CxAPI.ServerIP = CxAPI.CxConfig["IP"].(string)
		CxAPI.CxUser = CxAPI.CxConfig["User"].(string)
		CxAPI.CxPassword = CxAPI.CxConfig["Password"].(string)
	} else {
		fmt.Println("请输入CxServer ip:") //192.168.40.230
		fmt.Scanln(&CxAPI.ServerIP)
		fmt.Println("请输入Cx登陆用户:")
		fmt.Scanln(&CxAPI.CxUser)
		fmt.Println("请输入Cx登陆密码:")
		fmt.Scanln(&CxAPI.CxPassword)

		//保存配置
		ConfigFile := config{CxAPI.ServerIP, CxAPI.CxUser, CxAPI.CxPassword}
		data, err := json.Marshal(ConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		//保存配置文件
		WriteFile(data)

	}

	//登陆
	CxAPI.status, _ = CxAPI.Login("login")
	if CxAPI.status == 200 {
		fmt.Println("登陆成功")
		//获取token
		CxAPI.status = CxAPI.GetCxToken("token")

	} else {
		fmt.Println("登陆失败，重新登陆")
	}

	if CxAPI.status == 200 {
		for true {
			fmt.Println("选择操作：1（默认设置扫描项目获取报告）、2（查看所有项目）、3（查看扫描项目详情）、4（查看所有TeamsID）、5（退出）")
			fmt.Scanln(&num)

			if num == "1" {
				var CxName, ZipPath, RegisterType string
				fmt.Println("请输入项目名称:")
				fmt.Scanln(&CxName)
				fmt.Println("请输入源代码zip包路径:")
				fmt.Scanln(&ZipPath) //BookStoreJava.zip
				fmt.Println("请输入生成报告格式（PDF/XML/CSV/RTF）:")
				fmt.Scanln(&RegisterType)

				CxAPI.RegisterName = CxName + "." + RegisterType
				//创建项目
				StatusCode, ID := CxAPI.CreateProject(CxName, "create_project_with_default_configuration")
				fmt.Println(StatusCode)
				if StatusCode == 201 {
					//上传zip源代码包
					CxAPI.ZipFilePath = ZipPath
					StatusCode = CxAPI.upload_source_code_zip_file(ID, "upload_source_code_zip_file")
					if StatusCode == 204 {
						fmt.Println("上传源代码zip包完成")
						//创建扫描
						var scanid string
						StatusCode, scanid = CxAPI.create_new_scan(ID, "create_new_scan")
						CxAPI.ScanID = scanid
						if StatusCode == 201 {
							fmt.Println("创建扫描成功")
							//获取队列状态
							for true {
								_, ScanStatus, ScanStatuId := CxAPI.get_all_scan_details_in_queue(scanid, "get_all_scan_details_in_queue")
								CxAPI.ScanStatus = ScanStatus
								fmt.Printf("队列状态：%s，状态id：%s\n", ScanStatus, ScanStatuId)
								//实现暂停5s
								time.Sleep(time.Duration(5) * time.Second)
								if ScanStatus == "Finished" {
									break;
								}
							}

						} else {
							fmt.Println("获取扫描状态失败")
						}
					} else {
						fmt.Println("上传源代码zip包失败")
					}
				} else {
					fmt.Println("创建项目失败")
				}

				//获取报告
				if CxAPI.ScanStatus == "Finished" {
					fmt.Println("获取报告中....login")
					time.Sleep(time.Duration(10) * time.Second)
					//注册报告
					statu, reportId := CxAPI.register_scan_report(CxAPI.ScanID, RegisterType, "register_scan_report")
					if statu == 202 {
						fmt.Println("注册报告成功")
						//生成扫描报告
						fmt.Println("生成报告中....login")
						time.Sleep(time.Duration(10) * time.Second)
						CxAPI.get_reports_by_id(reportId, "get_reports_by_id")
					}
				} else {
					fmt.Println("注册报告失败")
				}

			} else
			if num == "2" {
				var projectid string
				fmt.Println("查看所有项目请回车；如按照项目ID查询请输入ProjectId：")
				fmt.Scanln(&projectid)
				//获取所有项目的详情
				CxAPI.Getprojects(projectid, "get_all_project_details")
			} else
			if num == "3" {
				var scanid string
				fmt.Println("查看所有扫描项目请回车；如按照扫描项目ID查询请输入ScanId：")
				fmt.Scanln(&scanid)
				//获取所有项目的详情
				CxAPI.get_sast_scan_details_by_scan_id(scanid, "get_sast_scan_details_by_scan_id")
			} else
			if num == "4" {
				CxAPI.GetTeams("get_all_teams")
			} else
			if num == "5" {
				fmt.Println("退出系统")
				break;
			} else {
				fmt.Println("输入错误请重新输入")
				continue;
			}

		}
	}

}
