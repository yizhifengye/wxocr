package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	imageName = "WX20250326-121647"
)

// OCRRequest 用于请求体中的图像字段
type OCRRequest struct {
	Image string `json:"image"`
}

// OCRResult 表示单个 OCR 识别结果
type OCRResult struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
	Rate   float64 `json:"rate"`
	Text   string  `json:"text"`
}

// OCRResponseResult 是外层 JSON 中的 "result" 字段的结构
type OCRResponseResult struct {
	ErrCode     int         `json:"errcode"`
	Height      int         `json:"height"`
	Width       int         `json:"width"`
	ImgPath     string      `json:"imgpath"`
	OCRResponse []OCRResult `json:"ocr_response"`
}

// OCRResponse 是整个 API 响应的结构
type OCRResponse struct {
	Result OCRResponseResult `json:"result"`
}

func ocrRecognize(imagePath string, imageUrl string, apiUrl string) (*OCRResponse, error) {
	var imgData []byte
	var err error

	// 读取本地文件或下载图像
	if imagePath != "" {
		imgData, err = os.ReadFile(imagePath)
		if err != nil {
			return nil, fmt.Errorf("读取本地图像失败: %v", err)
		}
	} else if imageUrl != "" {
		resp, err := http.Get(imageUrl)
		if err != nil {
			return nil, fmt.Errorf("下载图像失败: %v", err)
		}
		defer resp.Body.Close()
		imgData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取图像数据失败: %v", err)
		}
	} else {
		return nil, fmt.Errorf("请提供 imagePath 或 imageUrl")
	}

	// 转换为 base64
	base64Image := base64.StdEncoding.EncodeToString(imgData)
	return ocrRecognizeBase64(base64Image, apiUrl)
}

func ocrRecognizeBase64(base64Image string, apiUrl string) (*OCRResponse, error) {
	// 构造 JSON 请求
	reqBody, err := json.Marshal(OCRRequest{Image: base64Image})
	if err != nil {
		return nil, fmt.Errorf("JSON 编码失败: %v", err)
	}

	// 发送请求
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("API 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 API 响应失败: %v", err)
	}

	// 解析为 OCRResponse 结构体
	var ocrResponse OCRResponse
	err = json.Unmarshal(body, &ocrResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应 JSON 失败: %v", err)
	}

	// 打印结果
	fmt.Println("OCR识别结果:")
	for _, result := range ocrResponse.Result.OCRResponse {
		fmt.Printf("\"%s\" (%.0f,%.0f,%.0f,%.0f) %.1f%%\n", result.Text, result.Left, result.Right, result.Top, result.Bottom, result.Rate*100)
	}

	// 返回解析后的 OCRResponse
	return &ocrResponse, nil
}

func main() {
	apiUrl := "http://192.168.31.106:5000/ocr"
	imagePath := fmt.Sprintf("/Users/ace/Downloads/%s.png", imageName)
	_, err := ocrRecognize(imagePath, "", apiUrl)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}
}
