package ip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/tidwall/gjson"
	"net/http"
	"terraform-provider-fox/pkg/common"
	"time"
)

type IpInfo struct {
	Id           string    `json:"id"`
	Env          string    `json:"env"`
	GroupArr     []string  `json:"groupArr"`
	CidrArr      []string  `json:"cidrArr"`
	Created      time.Time `json:"created" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	LastModified time.Time `json:"lastModified" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

func convertByteToIpInfo(jsonByte []byte) (*IpInfo, error) {
	if jsonByte == nil {
		err := errors.New("convert fail: byte is null")
		return nil, err
	}
	ipInfo := IpInfo{}
	err := json.Unmarshal(jsonByte, &ipInfo)
	if err != nil {
		return nil, err
	}
	return &ipInfo, nil
}

func convertByteToIpInfoArr(jsonByte []byte) (*[]IpInfo, error) {
	if jsonByte == nil {
		err := errors.New("convert fail: byte is null")
		return nil, err
	}
	ipInfoArr := make([]IpInfo, 0)
	err := json.Unmarshal(jsonByte, &ipInfoArr)
	if err != nil {
		return nil, err
	}
	return &ipInfoArr, nil
}

func getHeaders(contentType string) map[string]string {
	h := make(map[string]string)
	if contentType != "" {
		h["Content-Type"] = contentType
	}
	return h
}

func GetIpInfo(config common.Config, id string, diagsPtr *diag.Diagnostics) (*IpInfo, error) {
	// request api
	apiPath := fmt.Sprintf("ip-envs/%s", id)
	headers := getHeaders("")
	status, respBody, err := common.HttpRequest(&common.HttpRequestArgs{
		Method:     http.MethodGet,
		Url:        common.GetApiUrl(config.Address, apiPath),
		TimeoutSec: 10,
		Headers:    headers,
	})
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "request api fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	} else if status != http.StatusOK {
		err := errors.New(fmt.Sprintf("get status: %d", status))
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "http response not ok",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	}

	// convert json
	ipInfoStr := gjson.Get(string(respBody), "res.ipInfo").String()
	ipInfo, err := convertByteToIpInfo([]byte(ipInfoStr))
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "convert json to IpInfo fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	}

	return ipInfo, nil
}

func GetIpInfoGroup(config common.Config, gId string, diagsPtr *diag.Diagnostics) (*[]IpInfo, error) {
	// request api
	apiPath := fmt.Sprintf("ip-groups/%s", gId)
	headers := getHeaders("")
	status, respBody, err := common.HttpRequest(&common.HttpRequestArgs{
		Method:     http.MethodGet,
		Url:        common.GetApiUrl(config.Address, apiPath),
		TimeoutSec: 10,
		Headers:    headers,
	})
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "request api fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	} else if status != http.StatusOK {
		err := errors.New(fmt.Sprintf("get status: %d", status))
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "http response not ok",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	}

	// convert json
	ipInfoArrStr := gjson.Get(string(respBody), "res.ipInfoArr").String()
	ipInfoArr, err := convertByteToIpInfoArr([]byte(ipInfoArrStr))
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "convert json to ipInfoArr fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	}

	return ipInfoArr, nil
}

func GetIpInfoAll(config common.Config, diagsPtr *diag.Diagnostics) (*[]IpInfo, error) {
	// request api
	apiPath := fmt.Sprintf("ip-all")
	headers := getHeaders("")
	status, respBody, err := common.HttpRequest(&common.HttpRequestArgs{
		Method:     http.MethodGet,
		Url:        common.GetApiUrl(config.Address, apiPath),
		TimeoutSec: 10,
		Headers:    headers,
	})
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "request api fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	} else if status != http.StatusOK {
		err := errors.New(fmt.Sprintf("get status: %d", status))
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "http response not ok",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	}

	// convert json
	ipInfoArrStr := gjson.Get(string(respBody), "res.ipInfoArr").String()
	ipInfoArr, err := convertByteToIpInfoArr([]byte(ipInfoArrStr))
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "convert json to ipInfoArr fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return nil, err
	}

	return ipInfoArr, nil
}

func CreateIpInfo(config common.Config, env string, groups []interface{}, cidrs []interface{}, diagsPtr *diag.Diagnostics) error {
	infoMap := make(map[string]interface{})
	infoMap["groupArr"] = groups
	infoMap["cidrArr"] = cidrs
	infoMap["managedBy"] = "tfc"

	// request api
	apiPath := fmt.Sprintf("ip-envs/%s", env)
	headers := getHeaders("application/json")
	reqBody, _ := json.Marshal(infoMap)
	status, _, err := common.HttpRequest(&common.HttpRequestArgs{
		Method:     http.MethodPost,
		Url:        common.GetApiUrl(config.Address, apiPath),
		TimeoutSec: 10,
		Body:       bytes.NewBuffer(reqBody),
		Headers:    headers,
	})
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "request api fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return err
	} else if status != http.StatusOK {
		err := errors.New(fmt.Sprintf("get status: %d", status))
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "http response not ok",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return err
	}

	return nil
}

func UpdateIpInfo(config common.Config, env string, groups []interface{}, cidrs []interface{}, diagsPtr *diag.Diagnostics) error {
	infoMap := make(map[string]interface{})
	infoMap["groupArr"] = groups
	infoMap["cidrArr"] = cidrs

	// request api
	apiPath := fmt.Sprintf("ip-envs/%s", env)
	headers := getHeaders("application/json")
	reqBody, _ := json.Marshal(infoMap)
	status, _, err := common.HttpRequest(&common.HttpRequestArgs{
		Method:     http.MethodPut,
		Url:        common.GetApiUrl(config.Address, apiPath),
		TimeoutSec: 10,
		Body:       bytes.NewBuffer(reqBody),
		Headers:    headers,
	})
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "request api fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return err
	} else if status != http.StatusOK {
		err := errors.New(fmt.Sprintf("get status: %d", status))
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "http response not ok",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return err
	}

	return nil
}

func DeleteIpInfo(config common.Config, env string, diagsPtr *diag.Diagnostics) error {
	// request api
	apiPath := fmt.Sprintf("ip-envs/%s", env)
	headers := getHeaders("")
	status, _, err := common.HttpRequest(&common.HttpRequestArgs{
		Method:     http.MethodDelete,
		Url:        common.GetApiUrl(config.Address, apiPath),
		TimeoutSec: 10,
		Headers:    headers,
	})
	if err != nil {
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "request api fail",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return err
	} else if status != http.StatusOK {
		err := errors.New(fmt.Sprintf("get status: %d", status))
		*diagsPtr = append(*diagsPtr, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "http response not ok",
			Detail:   fmt.Sprintf("error: %s", err.Error()),
		})
		return err
	}

	return nil
}
