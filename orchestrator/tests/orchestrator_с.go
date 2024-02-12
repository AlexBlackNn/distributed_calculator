package tests

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"orchestrator/internal/http-server/handlers/url/expression"
//	"time"
//)
//
////// TODO: get from config
////const (
////	host = "localhost:8080"
////)
////
////func TestExpression_HappyPath(t *testing.T) {
////	u := url.URL{
////		Scheme: "http",
////		Host:   host,
////	}
////	e := httpexpect.Default(t, u.String())
////
////	e.POST("/expression").
////		WithJSON(expression.Request{
////			Expression: "2+2*1",
////		}).
////		Expect().Status(http.StatusOK).JSON().Object().
////		ContainsKey("id").
////		ContainsKey("response")
////}
////
////func TestWholePathOfCalculation2_HappyPath(t *testing.T) {
////	var tests = []struct {
////		expression string
////		result     float64
////	}{
////		{"1+2", 3},
////		{"1+2*2", 5},
////		{"(1+2)*2", 6},
////		{"(1+2)*2/3", 2},
////		{"(1+2)*2/3*(1-2)", -2},
////	}
////
////	for _, test := range tests {
////		name := fmt.Sprintf("case(%d,%d)", test.expression, test.result)
////		t.Run(name, func(t *testing.T) {
////			got := IntMin(test.a, test.b)
////			if got != test.want {
////				t.Errorf("got %d, want %d", got, test.want)
////			}
////		})
////	}
////
////	u := url.URL{
////		Scheme: "http",
////		Host:   host,
////	}
////	e := httpexpect.Default(t, u.String())
////
////	obj := e.POST("/expression").
////		WithJSON(expression.Request{
////			Expression: "2+2*1",
////		}).
////		Expect().Status(http.StatusOK).JSON().Object().
////		ContainsKey("id").
////		ContainsKey("response")
////
////	var uidReturned string
////	obj.Value("id").Decode(&uidReturned)
////	fmt.Print(uidReturned)
////	time.Sleep(1 * time.Second)
////
////	result := e.GET("/expression/" + uidReturned).Expect().
////		Status(http.StatusOK).JSON().Object()
////	fmt.Println(result)
////
////	type ResponseBody struct {
////		Status string  `json:"status"`
////		Error  string  `json:"error,omitempty"`
////		Result float64 `json:"result,omitempty"`
////	}
////
////	resp := ResponseBody{}
////	result.Value("response").Decode(&resp)
////	fmt.Println(resp.Result, resp.Error, resp.Status)
////
////}
//
//type Interactor struct {
//	client *http.Client
//	url    string
//}
//
//func (ia *Interactor) PostExpression(
//	request expression.Request, response *expression.Response,
//) error {
//
//	jsonRequest, err := json.Marshal(request)
//	if err != nil {
//		return err
//	}
//	req, err := http.NewRequest("POST", ia.url, bytes.NewBuffer(jsonRequest))
//	req.Header.Set("accept-Type", "application/json")
//	req.Header.Set("Content-Type", "application/json")
//
//	resp, err := ia.client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//	err = json.Unmarshal(body, &response)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (ia *Interactor) GetResult(response *expression.Response) error {
//	req, err := http.NewRequest("GET", ia.url+"/"+response.Id, nil)
//	resp, err := ia.client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	fmt.Println("response Status:", resp.Status)
//	fmt.Println("response Headers:", resp.Header)
//	body, _ := io.ReadAll(resp.Body)
//	err = json.Unmarshal(body, &response)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func main() {
//	interactor := Interactor{
//		client: &http.Client{Timeout: 3 * time.Second},
//		url:    "http://localhost:8080/expression",
//	}
//
//	request := expression.Request{Expression: "1*10"}
//	response := expression.Response{}
//	err := interactor.PostExpression(request, &response)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(response)
//	err = interactor.GetResult(&response)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(response)
//}
//
////func TestExpression_HappyPath(t *testing.T) {
////	u := url.URL{
////		Scheme: "http",
////		Host:   host,
////	}
////	e := httpexpect.Default(t, u.String())
////
////	obj := e.POST("/expression").
////		WithJSON(expression.Request{
////			Expression: "2+2*1",
////		}).
////		Expect().Status(http.StatusOK).JSON().Object().
////		ContainsKey("id").
////		ContainsKey("response")
////
////	var uidReturned string
////	obj.Value("id").Decode(&uidReturned)
////	fmt.Print(uidReturned)
////	time.Sleep(1 * time.Second)
////
////	result := e.GET("/expression/" + uidReturned).Expect().
////		Status(http.StatusOK).JSON().Object()
////	fmt.Println(result)
////
////	type ResponseBody struct {
////		Status string  `json:"status"`
////		Error  string  `json:"error,omitempty"`
////		Result float64 `json:"result,omitempty"`
////	}
////
////	resp := ResponseBody{}
////	result.Value("response").Decode(&resp)
////	fmt.Println(resp.Result, resp.Error, resp.Status)
////
////}
