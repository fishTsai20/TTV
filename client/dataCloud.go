package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	ApiKey  string
	BaseURL string
}
type QueryData struct {
	QueryParameters map[string]string `json:"queryParameters"`
}

type ExecutionResponse struct {
	Data []struct {
		ExecutionId string `json:"executionId"`
		Status      string `json:"status"`
		Progress    int    `json:"progress"`
	} `json:"data"`
}

type Response struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

// Data structure for the "data" field
type Data struct {
	Columns       []Column        `json:"columns"`
	Data          [][]interface{} `json:"data"` // Assuming `data` will always contain integers
	ExecutionID   string          `json:"execution_id"`
	Message       string          `json:"message"`
	QueryID       int             `json:"query_id"`
	Status        string          `json:"status"`
	TotalRowCount int             `json:"total_row_count"`
}

// Column structure for the "columns" field
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (c *Client) executeQuery(queryId string, param map[string]string) (string, error) {
	url := fmt.Sprintf("%s/query/%s/execute", c.BaseURL, queryId)

	var (
		req *http.Request
		err error
	)

	if param != nil {
		queryData := QueryData{}
		queryData.QueryParameters = param

		data, err := json.Marshal(queryData)
		if err != nil {
			return "", err
		}

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			return "", err
		}
	} else {
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			return "", err
		}
	}

	req.Header.Set("X-API-KEY", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var execResponse ExecutionResponse
	if err := json.Unmarshal(body, &execResponse); err != nil {
		return "", err
	}

	return execResponse.Data[0].ExecutionId, nil
}

func (c *Client) checkStatus(executionId string) (ExecutionResponse, error) {
	url := fmt.Sprintf("%s/execution/%s/status", c.BaseURL, executionId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ExecutionResponse{}, err
	}

	req.Header.Set("X-API-KEY", c.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ExecutionResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ExecutionResponse{}, err
	}

	var statusResponse ExecutionResponse
	if err := json.Unmarshal(body, &statusResponse); err != nil {
		return ExecutionResponse{}, err
	}

	return statusResponse, nil
}

func (c *Client) getResults(executionId string) (*Response, error) {
	url := fmt.Sprintf("%s/execution/%s/results", c.BaseURL, executionId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", c.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret Response
	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Client) Query(queryId string, params map[string]string, msg api.Chattable, send func(msg api.Chattable) api.Message) (*Response, error) {
	executionId, err := c.executeQuery(queryId, params)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	var messageId int
	var status string
	for status != "FINISHED" && status != "FAILED" {
		statusResponse, err := c.checkStatus(executionId)
		if err != nil {
			fmt.Println("Error checking status:", err)
			return nil, err
		}

		status = statusResponse.Data[0].Status
		progress := statusResponse.Data[0].Progress
		totalSteps := 10
		curStep := int(math.Floor(float64(progress) / 100 * float64(totalSteps)))
		text := fmt.Sprintf("Progress: [%s%s] %d%%",
			strings.Repeat("■", curStep),
			strings.Repeat("□", totalSteps-curStep),
			progress)
		switch v := msg.(type) {
		case api.MessageConfig:
			if messageId == 0 {
				v.Text = text
				message := send(v)
				messageId = message.MessageID
			} else {
				msgs := api.NewEditMessageText(v.ChatID, messageId, text)
				send(msgs)
			}

		case api.EditMessageTextConfig:
			v.Text = text
			send(v)
		}

		time.Sleep(1 * time.Second)
	}

	return c.getResults(executionId)
}
