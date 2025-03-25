# WeChat OCR API Docker

基于微信OCR引擎的Docker化REST API服务。

[English](README_EN.MD) | 中文

## 概述

本项目将优秀的[wechat-ocr](https://github.com/swigger/wechat-ocr)项目中的微信OCR功能包装成一个简单的REST API服务，可以使用Docker轻松部署。它允许您利用微信强大的OCR功能对图像进行文字识别。

## ⚠️ 警告

**这是一个仅用于学习和交流目的的开源项目。请勿将其用于商业活动。用户对因不当使用本项目而导致的任何后果负有全部责任。**

**版权声明：** 本项目仅是对现有开源项目的容器化。如果您认为本仓库侵犯了您的版权或知识产权，请立即联系仓库所有者，仓库将会被迅速删除。我们尊重知识产权，无意侵犯任何权利。

## 致谢

本项目的实现离不开[swigger](https://github.com/swigger)及其[wechat-ocr](https://github.com/swigger/wechat-ocr)项目的工作。他们在逆向工程和创建微信OCR功能可用接口方面的努力构成了本服务的基础。

## 快速开始

### 使用Docker

```bash
# 拉取镜像
docker pull golangboyme/wxocr

# 运行容器
docker run -d -p 5000:5000 --name wechat-ocr-api golangboyme/wxocr
```

### API使用

向`/ocr`发送POST请求，JSON负载中包含您的base64编码图像：

```bash
curl -X POST http://localhost:5000/ocr \
  -H "Content-Type: application/json" \
  -d '{"image": "BASE64编码的图像数据"}'
```

#### 示例响应

```json
{
  "errcode": 0,
  "height": 72,
  "width": 410,
  "imgpath": "temp/5726fe7b-25d6-43a6-a50d-35b5f668fbb6.png",
  "ocr_response": [
    {
      "text": "aacss",
      "left": 80.63632202148438,
      "top": 29.634929656982422,
      "right": 236.47093200683594,
      "bottom": 55.28932189941406,
      "rate": 0.9997046589851379
    },
    {
      "text": "xxzsa",
      "left": 312.625,
      "top": 30.75,
      "right": 395.265625,
      "bottom": 55.09375,
      "rate": 0.997739315032959
    }
  ]
}
```

### Python客户端示例

这是一个简单的Python客户端，用于使用OCR API：

```python
import requests
import base64
import os

def ocr_recognize(image_path=None, image_url=None, api_url="http://localhost:5000/ocr"):
    """
    将图像发送到OCR API服务并获取识别结果。
    使用image_path或image_url（需要提供其中一个）。
    """
    # 获取图像数据
    if image_path:
        if not os.path.exists(image_path):
            print(f"错误：未找到本地图像：{image_path}")
            return
        with open(image_path, "rb") as image_file:
            img_data = image_file.read()
    elif image_url:
        try:
            response = requests.get(image_url)
            response.raise_for_status()
            img_data = response.content
        except Exception as e:
            print(f"下载图像失败：{str(e)}")
            return
    else:
        print("请提供image_path或image_url")
        return

    # 将图像转换为base64
    base64_image = base64.b64encode(img_data).decode('utf-8')
    
    # 发送请求到API
    try:
        response = requests.post(api_url, json={"image": base64_image})
        response.raise_for_status()
        return response.json()
    except Exception as e:
        print(f"API请求失败：{str(e)}")
        return None

# 使用示例
if __name__ == "__main__":
    # 本地图像示例
    result = ocr_recognize(image_path="ocrtest.png")
    if result:
        print(result)
    
    # URL图像示例（取消注释以使用）
    # result = ocr_recognize(image_url="https://example.com/image.png")
```

## 项目结构

- `main.py`：处理OCR请求的Flask API服务
- `opt/wechat/wxocr`：微信OCR二进制文件
- `opt/wechat/`：微信运行时依赖

## 技术细节

该服务使用Flask应用程序为微信OCR功能提供REST API接口。当提交图像时：

1. 解码base64编码的图像
2. 创建临时文件
3. 通过wcocr Python绑定由微信OCR引擎处理图像
4. 以JSON格式返回结果
5. 清理临时文件

## 限制

- 目前仅支持PNG图像（如需要可以扩展）
- 依赖于可能被微信更新的微信OCR二进制文件

## 许可证

本项目根据MIT许可证提供，并附加以下条款：

1. 仓库所有者对本软件的任何使用不承担任何责任。
2. 通过使用、复制、修改或分发本软件，您承认使用本软件的任何行为完全由您自行负责。
3. 对本仓库及其内容的任何使用均由用户自行承担风险和责任。
4. 仓库所有者对使用本软件产生的任何后果（法律上的或其他方面的）不承担责任。

## 贡献

欢迎贡献！请随时提交Pull Request。 