package main

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
	"reflect"
	"strings"
)

type projects struct {
	Id   int
	Name string
	Link interface{}
}
type details struct {
	Stage string
	Step  string
}
type status struct {
	Id       int
	Name     string
	Detailss details
}
type alljson struct {
	Id      int
	Project projects
	Status  status
}

//获取cookies
func Authentication() {
	Cliet := &http.Client{}
	request, err := http.NewRequest("POST", "http://192.168.40.180/cxrestapi/auth/login", strings.NewReader("userName=admin&password=Password01."))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//request.Header.Set("Cookie", "name=anny")
	response, err := Cliet.Do(request)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(body)
	Cxcookies := response.Cookies()
	fmt.Println(Cxcookies[0])
	fmt.Println(response.StatusCode)

}

//获取token
func GetCxToken(URL string) (token string) {
	Cliet := &http.Client{}
	request, err := http.NewRequest("POST", URL+"cxrestapi/auth/identity/connect/token", strings.NewReader("userName=admin&password=Password01.&grant_type=password&scope=sast_rest_api&client_id=resource_owner_client&client_secret=014DF517-39D1-4453-B7B3-9930C563627C"))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//request.Header.Set("Cookie", "name=anny")
	response, err := Cliet.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	//CxTokens := string(body)
	//在使用interface表示任何类型时，如果要将interface转为某一类型，直接强制转换是不行的，需要进行type assertion类型断言。
	//var a interface{} =access_token
	//token =a.(string)
	var Cxtoken map[string]interface{}
	json.Unmarshal(body, &Cxtoken)
	access_token := Cxtoken["access_token"]
	token = access_token.(string)
	//fmt.Println(reflect.TypeOf(access_token))
	//fmt.Println(reflect.TypeOf(token))
	//fmt.Println(access_token)
	//fmt.Println(response.StatusCode)
	return token
}

//登陆
func Login() {
	client := &http.Client{}
	request, err := http.NewRequest("POST", "http://192.168.40.216/cxrestapi/auth/login", strings.NewReader("userName=admin&password=Password01."))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, _ := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.StatusCode, body)

}

//获取所有项目数
func GetProjects() {
	Client := &http.Client{}
	request, err := http.NewRequest("GET", "http://192.168.40.216/cxrestapi/projects", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json;v=2.0")
	request.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSIsImtpZCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSJ9.eyJpc3MiOiJodHRwOi8vV0lOLUdPRlA1SThKOURCL2N4cmVzdGFwaS9hdXRoL2lkZW50aXR5IiwiYXVkIjoiaHR0cDovL1dJTi1HT0ZQNUk4SjlEQi9jeHJlc3RhcGkvYXV0aC9pZGVudGl0eS9yZXNvdXJjZXMiLCJleHAiOjE1NTcxOTQwOTgsIm5iZiI6MTU1NzEwNzY5OCwiY2xpZW50X2lkIjoicmVzb3VyY2Vfb3duZXJfY2xpZW50Iiwic2NvcGUiOiJzYXN0X3Jlc3RfYXBpIiwic3ViIjoiMiIsImF1dGhfdGltZSI6MTU1NzEwNzY5OCwiaWRwIjoiaWRzcnYiLCJpZCI6IjIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhZG1pbiIsImdpdmVuX25hbWUiOiJhZG1pbiIsImZhbWlseV9uYW1lIjoiYWRtaW4iLCJuYW1lIjoiYWRtaW4gYWRtaW4iLCJMQ0lEIjoiMjA1MiIsImVtYWlsIjoiYWRtaW5AY3guY29tIiwiVGVhbSI6IlxcMDAwMDAwMDAtMTExMS0xMTExLWIxMTEtOTg5YzkwNzBlYjExIiwic2FzdF9yb2xlIjpbInNhdmUtc2FzdC1zY2FuIiwiZGVsZXRlLXNhc3Qtc2NhbiIsInNhdmUtb3NhLXNjYW4iLCJtYW5hZ2UtcmVzdWx0cyIsIm1hbmFnZS1yZXN1bHQtY29tbWVudCIsIm1hbmFnZS1yZXN1bHQtZXhwbG9pdGFiaWxpdHkiLCJtYW5hZ2UtcmVzdWx0LXNldmVyaXR5IiwibWFuYWdlLWRhdGEtYW5hbHlzaXMtdGVtcGxhdGVzIiwidmlldy1kYXNoYm9hcmQiLCJnZW5lcmF0ZS1zY2FuLXJlcG9ydCIsIm1hbmFnZS1xdWVyaWVzIiwib3Blbi1pc3N1ZS10cmFja2luZy10aWNrZXRzIiwibWFuYWdlLWF1dGhlbnRpY2F0aW9uLXByb3RvY29scyIsIm1hbmFnZS1kYXRhLXJldGVudGlvbiIsIm1hbmFnZS1lbmdpbmUtc2VydmVycyIsIm1hbmFnZS1zeXN0ZW0tc2V0dGluZ3MiLCJ1c2Utb2RhdGEiLCJtYW5hZ2UtZXh0ZXJuYWwtc2VydmljZXMtc2V0dGluZ3MiLCJtYW5hZ2UtY3VzdG9tLWRlc2NyaXB0aW9uIiwidmlldy1hcHBzZWMtY29hY2gtc3RhdGlzdGljcyIsIm1hbmFnZS1jdXN0b20tZmllbGRzIiwic2F2ZS1wcm9qZWN0IiwibWFuYWdlLXVzZXJzIiwibWFuYWdlLWlzc3VlLXRyYWNraW5nLXN5c3RlbXMiLCJkZWxldGUtcHJvamVjdCIsIm1hbmFnZS1wcmUtcG9zdC1zY2FuLWFjdGlvbnMiLCJ2aWV3LWZhaWxlZC1zYXN0LXNjYW4iLCJjcmVhdGUtcHJlc2V0IiwidXBkYXRlLWFuZC1kZWxldGUtcHJlc2V0Il0sImFtciI6WyJwYXNzd29yZCJdfQ.R1k8rE5gZuw_7Bo-If9eq8yAaEKNLAJ8ms-XJLvRqXOtocXZB28UIo3jBsbupLjR5_LIH0h2ZhQCch9v_XRrD2hj-Cg0nSIlxHwsJF-8o0ZVyNEKW_gVVNrU8M3WfgheFJemJ125RUgg0leI929NE6LAoZiS3J3cNvlRNXU1lrXIbNT89WP3pWSFwYbGANx8bsb5m6VLMj7cgHiO1oNrlEP3_Gh80LcD913OrEylfZ6VRvehVwaRU_4Czi7etlxmvXZPn9PR_wVX5j-DyVdob9IM9DyTsFo3A4L5CPmxi-eGYGF5XWKDuIYax_esvmgaY_MEpzeCwA__T21qwjouqA")
	response, err := Client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

}

//按照ID查看项目
func GetIDproject(ID string) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://192.168.40.216/cxrestapi/projects/"+ID, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json;v=2.0")
	request.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSIsImtpZCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSJ9.eyJpc3MiOiJodHRwOi8vV0lOLUdPRlA1SThKOURCL2N4cmVzdGFwaS9hdXRoL2lkZW50aXR5IiwiYXVkIjoiaHR0cDovL1dJTi1HT0ZQNUk4SjlEQi9jeHJlc3RhcGkvYXV0aC9pZGVudGl0eS9yZXNvdXJjZXMiLCJleHAiOjE1NTcxOTQwOTgsIm5iZiI6MTU1NzEwNzY5OCwiY2xpZW50X2lkIjoicmVzb3VyY2Vfb3duZXJfY2xpZW50Iiwic2NvcGUiOiJzYXN0X3Jlc3RfYXBpIiwic3ViIjoiMiIsImF1dGhfdGltZSI6MTU1NzEwNzY5OCwiaWRwIjoiaWRzcnYiLCJpZCI6IjIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhZG1pbiIsImdpdmVuX25hbWUiOiJhZG1pbiIsImZhbWlseV9uYW1lIjoiYWRtaW4iLCJuYW1lIjoiYWRtaW4gYWRtaW4iLCJMQ0lEIjoiMjA1MiIsImVtYWlsIjoiYWRtaW5AY3guY29tIiwiVGVhbSI6IlxcMDAwMDAwMDAtMTExMS0xMTExLWIxMTEtOTg5YzkwNzBlYjExIiwic2FzdF9yb2xlIjpbInNhdmUtc2FzdC1zY2FuIiwiZGVsZXRlLXNhc3Qtc2NhbiIsInNhdmUtb3NhLXNjYW4iLCJtYW5hZ2UtcmVzdWx0cyIsIm1hbmFnZS1yZXN1bHQtY29tbWVudCIsIm1hbmFnZS1yZXN1bHQtZXhwbG9pdGFiaWxpdHkiLCJtYW5hZ2UtcmVzdWx0LXNldmVyaXR5IiwibWFuYWdlLWRhdGEtYW5hbHlzaXMtdGVtcGxhdGVzIiwidmlldy1kYXNoYm9hcmQiLCJnZW5lcmF0ZS1zY2FuLXJlcG9ydCIsIm1hbmFnZS1xdWVyaWVzIiwib3Blbi1pc3N1ZS10cmFja2luZy10aWNrZXRzIiwibWFuYWdlLWF1dGhlbnRpY2F0aW9uLXByb3RvY29scyIsIm1hbmFnZS1kYXRhLXJldGVudGlvbiIsIm1hbmFnZS1lbmdpbmUtc2VydmVycyIsIm1hbmFnZS1zeXN0ZW0tc2V0dGluZ3MiLCJ1c2Utb2RhdGEiLCJtYW5hZ2UtZXh0ZXJuYWwtc2VydmljZXMtc2V0dGluZ3MiLCJtYW5hZ2UtY3VzdG9tLWRlc2NyaXB0aW9uIiwidmlldy1hcHBzZWMtY29hY2gtc3RhdGlzdGljcyIsIm1hbmFnZS1jdXN0b20tZmllbGRzIiwic2F2ZS1wcm9qZWN0IiwibWFuYWdlLXVzZXJzIiwibWFuYWdlLWlzc3VlLXRyYWNraW5nLXN5c3RlbXMiLCJkZWxldGUtcHJvamVjdCIsIm1hbmFnZS1wcmUtcG9zdC1zY2FuLWFjdGlvbnMiLCJ2aWV3LWZhaWxlZC1zYXN0LXNjYW4iLCJjcmVhdGUtcHJlc2V0IiwidXBkYXRlLWFuZC1kZWxldGUtcHJlc2V0Il0sImFtciI6WyJwYXNzd29yZCJdfQ.R1k8rE5gZuw_7Bo-If9eq8yAaEKNLAJ8ms-XJLvRqXOtocXZB28UIo3jBsbupLjR5_LIH0h2ZhQCch9v_XRrD2hj-Cg0nSIlxHwsJF-8o0ZVyNEKW_gVVNrU8M3WfgheFJemJ125RUgg0leI929NE6LAoZiS3J3cNvlRNXU1lrXIbNT89WP3pWSFwYbGANx8bsb5m6VLMj7cgHiO1oNrlEP3_Gh80LcD913OrEylfZ6VRvehVwaRU_4Czi7etlxmvXZPn9PR_wVX5j-DyVdob9IM9DyTsFo3A4L5CPmxi-eGYGF5XWKDuIYax_esvmgaY_MEpzeCwA__T21qwjouqA")
	/*设置URL后面？的参数
	q := request.URL.Query()
	q.Add("123","123")
	request.URL.RawQuery = q.Encode()
	*/
	//查看url
	fmt.Println(request.URL.String())
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println(response.StatusCode)
}

//按照ID更新项目
func GetIDupdate(ID string) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", "http://192.168.40.216/cxrestapi/projects/"+ID, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json;v=2.0")
	request.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSIsImtpZCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSJ9.eyJpc3MiOiJodHRwOi8vV0lOLUdPRlA1SThKOURCL2N4cmVzdGFwaS9hdXRoL2lkZW50aXR5IiwiYXVkIjoiaHR0cDovL1dJTi1HT0ZQNUk4SjlEQi9jeHJlc3RhcGkvYXV0aC9pZGVudGl0eS9yZXNvdXJjZXMiLCJleHAiOjE1NTcxOTQwOTgsIm5iZiI6MTU1NzEwNzY5OCwiY2xpZW50X2lkIjoicmVzb3VyY2Vfb3duZXJfY2xpZW50Iiwic2NvcGUiOiJzYXN0X3Jlc3RfYXBpIiwic3ViIjoiMiIsImF1dGhfdGltZSI6MTU1NzEwNzY5OCwiaWRwIjoiaWRzcnYiLCJpZCI6IjIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhZG1pbiIsImdpdmVuX25hbWUiOiJhZG1pbiIsImZhbWlseV9uYW1lIjoiYWRtaW4iLCJuYW1lIjoiYWRtaW4gYWRtaW4iLCJMQ0lEIjoiMjA1MiIsImVtYWlsIjoiYWRtaW5AY3guY29tIiwiVGVhbSI6IlxcMDAwMDAwMDAtMTExMS0xMTExLWIxMTEtOTg5YzkwNzBlYjExIiwic2FzdF9yb2xlIjpbInNhdmUtc2FzdC1zY2FuIiwiZGVsZXRlLXNhc3Qtc2NhbiIsInNhdmUtb3NhLXNjYW4iLCJtYW5hZ2UtcmVzdWx0cyIsIm1hbmFnZS1yZXN1bHQtY29tbWVudCIsIm1hbmFnZS1yZXN1bHQtZXhwbG9pdGFiaWxpdHkiLCJtYW5hZ2UtcmVzdWx0LXNldmVyaXR5IiwibWFuYWdlLWRhdGEtYW5hbHlzaXMtdGVtcGxhdGVzIiwidmlldy1kYXNoYm9hcmQiLCJnZW5lcmF0ZS1zY2FuLXJlcG9ydCIsIm1hbmFnZS1xdWVyaWVzIiwib3Blbi1pc3N1ZS10cmFja2luZy10aWNrZXRzIiwibWFuYWdlLWF1dGhlbnRpY2F0aW9uLXByb3RvY29scyIsIm1hbmFnZS1kYXRhLXJldGVudGlvbiIsIm1hbmFnZS1lbmdpbmUtc2VydmVycyIsIm1hbmFnZS1zeXN0ZW0tc2V0dGluZ3MiLCJ1c2Utb2RhdGEiLCJtYW5hZ2UtZXh0ZXJuYWwtc2VydmljZXMtc2V0dGluZ3MiLCJtYW5hZ2UtY3VzdG9tLWRlc2NyaXB0aW9uIiwidmlldy1hcHBzZWMtY29hY2gtc3RhdGlzdGljcyIsIm1hbmFnZS1jdXN0b20tZmllbGRzIiwic2F2ZS1wcm9qZWN0IiwibWFuYWdlLXVzZXJzIiwibWFuYWdlLWlzc3VlLXRyYWNraW5nLXN5c3RlbXMiLCJkZWxldGUtcHJvamVjdCIsIm1hbmFnZS1wcmUtcG9zdC1zY2FuLWFjdGlvbnMiLCJ2aWV3LWZhaWxlZC1zYXN0LXNjYW4iLCJjcmVhdGUtcHJlc2V0IiwidXBkYXRlLWFuZC1kZWxldGUtcHJlc2V0Il0sImFtciI6WyJwYXNzd29yZCJdfQ.R1k8rE5gZuw_7Bo-If9eq8yAaEKNLAJ8ms-XJLvRqXOtocXZB28UIo3jBsbupLjR5_LIH0h2ZhQCch9v_XRrD2hj-Cg0nSIlxHwsJF-8o0ZVyNEKW_gVVNrU8M3WfgheFJemJ125RUgg0leI929NE6LAoZiS3J3cNvlRNXU1lrXIbNT89WP3pWSFwYbGANx8bsb5m6VLMj7cgHiO1oNrlEP3_Gh80LcD913OrEylfZ6VRvehVwaRU_4Czi7etlxmvXZPn9PR_wVX5j-DyVdob9IM9DyTsFo3A4L5CPmxi-eGYGF5XWKDuIYax_esvmgaY_MEpzeCwA__T21qwjouqA")
	/*设置URL后面？的参数
	q := request.URL.Query()
	q.Add("123","123")
	request.URL.RawQuery = q.Encode()
	*/
	//查看url
	fmt.Println(request.URL.String())
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println(response.StatusCode)
}

//按照项目ID删除项目
func GetIDdelete(ID string) {
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", "http://192.168.40.216/cxrestapi/projects/"+ID, strings.NewReader("deleteRunningScans=true"))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json;v=2.0")
	request.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSIsImtpZCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSJ9.eyJpc3MiOiJodHRwOi8vV0lOLUdPRlA1SThKOURCL2N4cmVzdGFwaS9hdXRoL2lkZW50aXR5IiwiYXVkIjoiaHR0cDovL1dJTi1HT0ZQNUk4SjlEQi9jeHJlc3RhcGkvYXV0aC9pZGVudGl0eS9yZXNvdXJjZXMiLCJleHAiOjE1NTcxOTQwOTgsIm5iZiI6MTU1NzEwNzY5OCwiY2xpZW50X2lkIjoicmVzb3VyY2Vfb3duZXJfY2xpZW50Iiwic2NvcGUiOiJzYXN0X3Jlc3RfYXBpIiwic3ViIjoiMiIsImF1dGhfdGltZSI6MTU1NzEwNzY5OCwiaWRwIjoiaWRzcnYiLCJpZCI6IjIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhZG1pbiIsImdpdmVuX25hbWUiOiJhZG1pbiIsImZhbWlseV9uYW1lIjoiYWRtaW4iLCJuYW1lIjoiYWRtaW4gYWRtaW4iLCJMQ0lEIjoiMjA1MiIsImVtYWlsIjoiYWRtaW5AY3guY29tIiwiVGVhbSI6IlxcMDAwMDAwMDAtMTExMS0xMTExLWIxMTEtOTg5YzkwNzBlYjExIiwic2FzdF9yb2xlIjpbInNhdmUtc2FzdC1zY2FuIiwiZGVsZXRlLXNhc3Qtc2NhbiIsInNhdmUtb3NhLXNjYW4iLCJtYW5hZ2UtcmVzdWx0cyIsIm1hbmFnZS1yZXN1bHQtY29tbWVudCIsIm1hbmFnZS1yZXN1bHQtZXhwbG9pdGFiaWxpdHkiLCJtYW5hZ2UtcmVzdWx0LXNldmVyaXR5IiwibWFuYWdlLWRhdGEtYW5hbHlzaXMtdGVtcGxhdGVzIiwidmlldy1kYXNoYm9hcmQiLCJnZW5lcmF0ZS1zY2FuLXJlcG9ydCIsIm1hbmFnZS1xdWVyaWVzIiwib3Blbi1pc3N1ZS10cmFja2luZy10aWNrZXRzIiwibWFuYWdlLWF1dGhlbnRpY2F0aW9uLXByb3RvY29scyIsIm1hbmFnZS1kYXRhLXJldGVudGlvbiIsIm1hbmFnZS1lbmdpbmUtc2VydmVycyIsIm1hbmFnZS1zeXN0ZW0tc2V0dGluZ3MiLCJ1c2Utb2RhdGEiLCJtYW5hZ2UtZXh0ZXJuYWwtc2VydmljZXMtc2V0dGluZ3MiLCJtYW5hZ2UtY3VzdG9tLWRlc2NyaXB0aW9uIiwidmlldy1hcHBzZWMtY29hY2gtc3RhdGlzdGljcyIsIm1hbmFnZS1jdXN0b20tZmllbGRzIiwic2F2ZS1wcm9qZWN0IiwibWFuYWdlLXVzZXJzIiwibWFuYWdlLWlzc3VlLXRyYWNraW5nLXN5c3RlbXMiLCJkZWxldGUtcHJvamVjdCIsIm1hbmFnZS1wcmUtcG9zdC1zY2FuLWFjdGlvbnMiLCJ2aWV3LWZhaWxlZC1zYXN0LXNjYW4iLCJjcmVhdGUtcHJlc2V0IiwidXBkYXRlLWFuZC1kZWxldGUtcHJlc2V0Il0sImFtciI6WyJwYXNzd29yZCJdfQ.R1k8rE5gZuw_7Bo-If9eq8yAaEKNLAJ8ms-XJLvRqXOtocXZB28UIo3jBsbupLjR5_LIH0h2ZhQCch9v_XRrD2hj-Cg0nSIlxHwsJF-8o0ZVyNEKW_gVVNrU8M3WfgheFJemJ125RUgg0leI929NE6LAoZiS3J3cNvlRNXU1lrXIbNT89WP3pWSFwYbGANx8bsb5m6VLMj7cgHiO1oNrlEP3_Gh80LcD913OrEylfZ6VRvehVwaRU_4Czi7etlxmvXZPn9PR_wVX5j-DyVdob9IM9DyTsFo3A4L5CPmxi-eGYGF5XWKDuIYax_esvmgaY_MEpzeCwA__T21qwjouqA")
	/*设置URL后面？的参数
	q := request.URL.Query()
	q.Add("123","123")
	request.URL.RawQuery = q.Encode()
	*/
	//查看url
	fmt.Println(request.URL.String())
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println(response.StatusCode)
}

//获取teams的信息
func GetTeams() {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://192.168.40.225/cxrestapi/auth/teams", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSIsImtpZCI6IjZiMkxVZjRBSmg4NVhCdER4b1IzRE5mR2dMWSJ9.eyJpc3MiOiJodHRwOi8vV0lOLUdPRlA1SThKOURCL2N4cmVzdGFwaS9hdXRoL2lkZW50aXR5IiwiYXVkIjoiaHR0cDovL1dJTi1HT0ZQNUk4SjlEQi9jeHJlc3RhcGkvYXV0aC9pZGVudGl0eS9yZXNvdXJjZXMiLCJleHAiOjE1NTcyODI4NzQsIm5iZiI6MTU1NzE5NjQ3NCwiY2xpZW50X2lkIjoicmVzb3VyY2Vfb3duZXJfY2xpZW50Iiwic2NvcGUiOiJzYXN0X3Jlc3RfYXBpIiwic3ViIjoiMiIsImF1dGhfdGltZSI6MTU1NzE5NjQ3NCwiaWRwIjoiaWRzcnYiLCJpZCI6IjIiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhZG1pbiIsImdpdmVuX25hbWUiOiJhZG1pbiIsImZhbWlseV9uYW1lIjoiYWRtaW4iLCJuYW1lIjoiYWRtaW4gYWRtaW4iLCJMQ0lEIjoiMjA1MiIsImVtYWlsIjoiYWRtaW5AY3guY29tIiwiVGVhbSI6IlxcMDAwMDAwMDAtMTExMS0xMTExLWIxMTEtOTg5YzkwNzBlYjExIiwic2FzdF9yb2xlIjpbInNhdmUtc2FzdC1zY2FuIiwiZGVsZXRlLXNhc3Qtc2NhbiIsInNhdmUtb3NhLXNjYW4iLCJtYW5hZ2UtcmVzdWx0cyIsIm1hbmFnZS1yZXN1bHQtY29tbWVudCIsIm1hbmFnZS1yZXN1bHQtZXhwbG9pdGFiaWxpdHkiLCJtYW5hZ2UtcmVzdWx0LXNldmVyaXR5IiwibWFuYWdlLWRhdGEtYW5hbHlzaXMtdGVtcGxhdGVzIiwidmlldy1kYXNoYm9hcmQiLCJnZW5lcmF0ZS1zY2FuLXJlcG9ydCIsIm1hbmFnZS1xdWVyaWVzIiwib3Blbi1pc3N1ZS10cmFja2luZy10aWNrZXRzIiwibWFuYWdlLWF1dGhlbnRpY2F0aW9uLXByb3RvY29scyIsIm1hbmFnZS1kYXRhLXJldGVudGlvbiIsIm1hbmFnZS1lbmdpbmUtc2VydmVycyIsIm1hbmFnZS1zeXN0ZW0tc2V0dGluZ3MiLCJ1c2Utb2RhdGEiLCJtYW5hZ2UtZXh0ZXJuYWwtc2VydmljZXMtc2V0dGluZ3MiLCJtYW5hZ2UtY3VzdG9tLWRlc2NyaXB0aW9uIiwidmlldy1hcHBzZWMtY29hY2gtc3RhdGlzdGljcyIsIm1hbmFnZS1jdXN0b20tZmllbGRzIiwic2F2ZS1wcm9qZWN0IiwibWFuYWdlLXVzZXJzIiwibWFuYWdlLWlzc3VlLXRyYWNraW5nLXN5c3RlbXMiLCJkZWxldGUtcHJvamVjdCIsIm1hbmFnZS1wcmUtcG9zdC1zY2FuLWFjdGlvbnMiLCJ2aWV3LWZhaWxlZC1zYXN0LXNjYW4iLCJjcmVhdGUtcHJlc2V0IiwidXBkYXRlLWFuZC1kZWxldGUtcHJlc2V0Il0sImFtciI6WyJwYXNzd29yZCJdfQ.Qh8M2a6Do9WJC_WZ34SWJbZFtoXlgwpnT31FfsoZ6FmhtDhk_rHhZ1iCWx1G1CSLisHu8RU-UFRo6LBj2cpPnwz2Wvl8IL5ufGNxAT-YXz45waFTDfo908fyYd6cOZmke4NdlHfUSuZ1FWVO8caX2daxNGviDg2fdl8D9I_eJT2wgWJcD-5fz875dC6xabTzbeA5QsukeeFD5z6VlMhV_tCVkAnLjPkmhvPijhQhUi10GbTRZWDPLpoDev8vS0hg_haO6v8IWHkBw68oKORNPMduphCQqZE1SZx6sjV5il1aYoTLK6eaki9kHPphoADBSJN1M7SjNWcIUk8I9MKK6A")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}

//使用默认配置创建项目
func CreateProject(CxToken string, URL string) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", URL+"cxrestapi/projects", strings.NewReader("name=API_test&owningTeam=00000000-1111-1111-b111-989c9070eb11&isPublic=true"))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", CxToken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	fmt.Println(response.StatusCode)
}

//上传源代码zip包
func UploadZip(CxToken, ID, URL string) {
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
	writer.Close() // 发送之前必须调用Close()以写入结尾行
	request, err := http.NewRequest("POST", URL+"cxrestapi/projects/"+ID+"/sourceCode/attachments", buf)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", CxToken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	fmt.Println(response.StatusCode)
}

//开始扫描
func ScanProject(CxToken, ID, URL string) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", URL+"/cxrestapi/sast/scans", strings.NewReader("projectId="+ID+"&isIncremental=false&isPublic=false&forceScan=True"))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", CxToken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	fmt.Println(response.StatusCode)
}

//查看所有扫描项目
func GetScanProjects(CxToken, URL string) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", URL+"cxrestapi/sast/scans", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", CxToken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	//var name map[string]interface{}
	//var a []  string
	var a []alljson
	er := json.Unmarshal(body, &a)
	if er != nil {
		log.Fatal(er)
	}
	//fmt.Println(string(body))
	num := len(a)
	//fmt.Println(string(body))
	fmt.Println(num)
	//fmt.Println(a[0])
	i := 0
	for i < num {
		fmt.Printf("扫描ID：%d   项目名称：%v   项目状态：%v\n", a[i].Id, a[i].Project.Name, a[i].Status.Name)
		i += 1
	}

}

//获取扫描报告
//1、注册扫描报告
func RegisterPeport(CxToken, URL, Filetype, scanID string) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", URL+"cxrestapi/reports/sastScan", strings.NewReader("reportType="+Filetype+"&scanId="+scanID))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", CxToken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	fmt.Println(response.StatusCode)
}

//2、获取生成报告的状态
func Getstatus(Cxtoken, URL, ID string) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", URL+"/cxrestapi/reports/sastScan/"+ID+"/status", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", Cxtoken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(response.StatusCode)
	fmt.Println(string(body))
}

//3、获取扫描报告
func Getreports(CxToken, URL, ID string) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", URL+"cxrestapi/reports/sastScan/"+ID, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json;v=1.0")
	request.Header.Set("Authorization", CxToken)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(response.StatusCode)
	fmt.Println(string(body))
	fmt.Println(reflect.TypeOf(body))
	//写入文件
	er := ioutil.WriteFile("test.XML", body, 0644)
	if er != nil {
		panic(er)
	}
}

func main() {
	URL := "http://192.168.40.134/"
	//Authentication()
	CxToken := "Bearer " + GetCxToken(URL)
	//fmt.Println(CxToken)
	//Login()
	//GetProjects()
	//GetIDproject("60057")
	//GetIDupdate("60057")
	//GetIDdelete("60057")
	//GetTeams()
	//CreateProject(CxToken, URL)
	//UploadZip(CxToken, "60060", URL)
	//ScanProject(CxToken, "60060", URL)
	//GetScanProjects(CxToken, URL)
	//RegisterPeport(CxToken, URL, "XML", "1140094")
	//Getstatus(CxToken, URL, "4040")
	Getreports(CxToken, URL, "4041")

}
