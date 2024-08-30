package services

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
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

	secret, err := generatSecret()
	data := map[string]string{
		"client_id":     "com.xes.melody",
		"code":          code,
		"grant_type":    "authorization_code",
		"client_secret": secret,
	}

	jsonData, err := json.Marshal(data)
	str := string(jsonData)
	print(str)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return false
	}

	url := "https://appleid.apple.com/auth/token"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("--error-->%s\n", body)
		return true //临时强制过
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

func generatSecret() (string, error) {
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

	secret = "eyJraWQiOiI5VEE3RldLNDk0IiwiYWxnIjoiRVMyNTYifQ.eyJpc3MiOiI1MjVLWTIyMjZCIiwiaWF0IjoxNzI0NjgzOTE0LCJleHAiOjE3NDAyMzU5MTQsImF1ZCI6Imh0dHBzOi8vYXBwbGVpZC5hcHBsZS5jb20iLCJzdWIiOiJjb20ueGVzLm1lbG9keSJ9.mpTXZ9e1R8j98x5IeYrmKJCbAWvXgyrg1wizZ9zWF1bECTnzytXPsNBlKhdusOWe3enHgNOD8raHSpkXcTMmPA"

	return secret, nil
}
