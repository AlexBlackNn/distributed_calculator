package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orchestrator/internal/http-server/handlers/calculation/expression"
)

type Interactor struct {
	client *http.Client
	url    string
}

func (ia *Interactor) PostExpression(
	request expression.Request, response *expression.Response,
) error {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", ia.url, bytes.NewBuffer(jsonRequest))
	req.Header.Set("accept-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := ia.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &response)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (ia *Interactor) GetResult(response *expression.Response) error {
	req, err := http.NewRequest("GET", ia.url+"/"+response.Id, nil)
	resp, err := ia.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	return nil
}
