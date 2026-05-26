package pkg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/bwmarrin/snowflake"
)

// 雪花ID
func GenerateID() string {
	node, _ := snowflake.NewNode(1)
	return node.Generate().String()
}

// MD5
func EncryptByMD5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// UploadToIPFS 上传文件到本地 IPFS 节点，返回 CID
func UploadToIPFS(file multipart.File, header *multipart.FileHeader) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return "", err
	}
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", "http://127.0.0.1:5001/api/v0/add", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	cid, ok := result["Hash"].(string)
	if !ok {
		return "", fmt.Errorf("IPFS 返回格式异常: %s", string(respBody))
	}
	return cid, nil
}
