package handle

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// Info ...
type Info struct {
	MsgText    string `json:"msgText"`    // 短信内容
	Destmobile string `json:"destmobile"` // 手机号码
}

const (
	messageServerURL = "http://www.jianzhou.sh.cn/JianzhouSMSWSServer/http/sendBatchMessage"
	account          = "sdk_aqjr"
	password         = "njmz504Aq808"
	signature        = "【国运产权】"
)

// Send 发生短信
func (i *Info) Send() (string, error) {
	resp, err := http.PostForm(messageServerURL,
		url.Values{
			"account":    {account},
			"password":   {password},
			"msgText":    {i.MsgText + signature},
			"destmobile": {i.Destmobile},
		},
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
