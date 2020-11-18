package handle

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/xiao0811/xiniu/config"
)

// Info ...
type Info struct {
	MsgText    string `json:"msgText"`    // 短信内容
	Destmobile string `json:"destmobile"` // 手机号码
}

var (
	conf             = config.Conf.MessageConfig
	messageServerURL = conf.ServerURL
	account          = conf.Account
	password         = conf.Password
	signature        = conf.Signature
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
