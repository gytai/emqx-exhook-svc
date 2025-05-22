package netbase

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type DeviceAuth struct {
	Owner          string `json:"owner"`
	OrgId          int64  `json:"orgId"`
	DeviceId       string `json:"deviceId"`
	DeviceType     string `json:"deviceType"`
	DeviceProtocol string `json:"deviceProtocol"`
	ProductId      string `json:"productId"`
	RuleChainId    string `json:"ruleChainId"`
	Name           string `json:"name"`
	CreatedAt      int64  `json:"created_at"`
	ExpiredAt      int64  `json:"expired_at"`
}

type DeviceEventInfo struct {
	DeviceId   string      `json:"deviceId"`
	DeviceAuth *DeviceAuth `json:"deviceAuth"`
	Datas      string      `json:"datas"`
	RequestId  string      `json:"requestId"`
	Topic      string      `json:"topic"`
}

func (j *DeviceEventInfo) Bytes() []byte {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return b
}

func (entity *DeviceAuth) GetDeviceToken(key string) error {
	if err := GetDeviceEtoken(key, entity); err != nil {
		return err
	}
	return nil
}

func (token *DeviceAuth) MD5ID() string {
	buf := bytes.NewBufferString(token.DeviceId)
	buf.WriteString(token.DeviceType)
	buf.WriteString(strconv.FormatInt(token.CreatedAt, 10))
	access := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
	access = strings.TrimRight(access, "=")
	return access
}

func (token *DeviceAuth) GetMarshal() string {
	marshal, _ := json.Marshal(*token)
	return string(marshal)
}

func (token *DeviceAuth) GetUnMarshal(data []byte) error {
	return json.Unmarshal(data, token)
}

// 序列化
func (m *DeviceAuth) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// 反序列化
func (m *DeviceAuth) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
