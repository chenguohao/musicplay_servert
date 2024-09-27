package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func jwtDecodeWithJwtString(tokenString string) (map[string]interface{}, error) {
	token, _ := jwt.Parse(tokenString, nil)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("unable to parse JWT")
}

func genSecret() string {
	// Generate your client_secret using your private key and other details
	// This is a placeholder for demonstration purposes
	return "your_client_secret"
}

func AppleSign(code, token string) bool {

	dict1, err := jwtDecodeWithJwtString(token)
	if err != nil {
		fmt.Println("Error decoding JWT:", err)
		return false
	}
	fmt.Println(">>解析原始的:", dict1)

	// 生成 client_secret
	secret, err := generatSecret()
	if err != nil {
		fmt.Println("Error generating client secret:", err)
		return false
	}

	// 创建表单数据
	data := url.Values{}
	data.Set("client_id", "com.xes.melody")
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", secret)

	// 将数据编码为 x-www-form-urlencoded 格式
	formData := data.Encode()

	// 打印请求体 (可选)
	fmt.Println("Form data:", formData)

	// 创建请求
	url := "https://appleid.apple.com/auth/token"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(formData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	//	dict1, err := jwtDecodeWithJwtString(token)
	//	if err != nil {
	//		fmt.Println("Error decoding JWT:", err)
	//		return false
	//	}
	//	fmt.Println(">>解析原始的:", dict1)
	//
	//	secret, err := generatSecret()
	//	data := map[string]string{
	//		"client_id":     "com.xes.melody",
	//		"code":          code,
	//		"grant_type":    "authorization_code",
	//		"client_secret": secret,
	//	}
	//
	//	jsonData, err := json.Marshal(data)
	//	str := string(jsonData)
	//	print(str)
	//	if err != nil {
	//		fmt.Println("Error marshalling JSON:", err)
	//		return false
	//	}
	//
	//	url := "https://appleid.apple.com/auth/token"
	//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	//	if err != nil {
	//		fmt.Println("Error creating request:", err)
	//		return false
	//	}
	//
	//	req.Header.Set("Content-Type", "application/json")
	//
	//	client := &http.Client{Timeout: 10 * time.Second}
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		fmt.Println("Error sending request:", err)
	//		return false
	//	}
	//	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("--error-->%s\n", body)
		return false //临时强制过
	}

	// Parsing the response into a map
	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Println("Error decoding response:", err)
		return false
	}

	fmt.Println("--success-->", tokenResponse)

	// Extracting and decoding the "id_token"
	idToken, ok := tokenResponse["id_token"].(string)
	if !ok {
		fmt.Println("Error: id_token not found or invalid")
		return false
	}

	dict2, err := jwtDecodeWithJwtString(idToken)
	if err != nil {
		fmt.Println("Error decoding JWT:", err)
		return false
	}
	fmt.Println(">>解析请求到的:", dict2)

	return true
}

//go:embed AuthKey_9TA7FWK494.p8
var privateKeyBytes []byte

func generatSecret() (string, error) {
	//keyFile := "services/AuthKey_9TA7FWK494.p8" // 从 Developer Center 后台下载的那个 p8 文件
	teamID := "525KY2226B"       // 开发者账号的 teamID
	clientID := "com.xes.melody" // 应用的 BundleID
	keyID := "9TA7FWK494"        // 从 Developer Center 后台找到 keyID
	validityPeriod := 180        // 有效期 180 天

	// 读取 .p8 文件中的私钥
	//privateKeyBytes, err := ioutil.ReadFile(keyFile)
	//if err != nil {
	//	fmt.Println("Error reading private key file:", err)
	//	return "", err
	//}

	// 解码 PEM 格式的私钥
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		fmt.Println("Failed to decode PEM block containing private key")
		return "", fmt.Errorf("failed to decode private key")
	}

	// 解析 PKCS8 格式的私钥
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing EC private key:", err)
		return "", err
	}

	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		fmt.Println("Private key is not of type ECDSA")
		return "", fmt.Errorf("private key is not ECDSA")
	}

	// 创建 JWT 令牌并设置 claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": teamID,
		"iat": time.Now().Unix(),                                                     // 当前时间戳
		"exp": time.Now().Add(time.Duration(validityPeriod) * 24 * time.Hour).Unix(), // 有效期
		"aud": "https://appleid.apple.com",
		"sub": clientID,
	})

	// 设置 header 中的 keyID
	token.Header["kid"] = keyID
	fmt.Println("token:", token)
	// 使用 ECDSA 私钥对令牌进行签名
	secret, err := token.SignedString(ecdsaPrivateKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	return secret, nil
}
func generatSecret2() (string, error) {
	keyFile := "services/AuthKey_9TA7FWK494.p8" // 从 Developer Center 后台下载的那个 p8 文件
	teamID := "525KY2226B"                      // 开发者账号的 teamID
	clientID := "com.xes.melody"                // 应用的 BundleID
	keyID := "9TA7FWK494"                       // 从 Developer Center 后台找到 keyID
	validityPeriod := 180                       // 有效期 180 天

	// 读取 .p8 文件中的私钥
	privateKeyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Error reading private key file:", err)
		return "", err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		fmt.Println("Failed to decode PEM block containing private key")
		return "", err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing EC private key:", err)
		return "", err
	}

	// 创建 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": teamID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(validityPeriod) * 24 * time.Hour).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": clientID,
	})

	token.Header["kid"] = keyID

	// 使用私钥对令牌进行签名
	secret, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	return secret, nil
}
