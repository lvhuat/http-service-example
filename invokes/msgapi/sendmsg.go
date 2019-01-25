package msgapi

import (
	"github.com/lworkltd/kits/service/invoke"
)

var (
	seviceAddr = "msgapi:8080"
)

type SendMailRequest struct {
	Email      string `json:"email"`
	TemplateId string `json:"templateId"`
	Parameters map[string]string
}

type SendMailResponse struct {
	DealId string `json:"dealId"`
}

func SendEmailToUser(req *SendMailRequest) (*SendMailResponse, error) {
	var dataRsp SendMailResponse
	err := invoke.Addr(seviceAddr).
		Post("/msgapi/v1/mail/send").
		Hystrix(100000, 0, 0).
		Json(req).
		Result(&dataRsp)
	if err != nil {
		return nil, err
	}
	return &dataRsp, nil
}

func SendRegisterVerifyMail(user, email, verifyCode string) (string, error) {
	req := &SendMailRequest{
		Email:      email,
		TemplateId: "REGISTER_VERIFY_CODE",
		Parameters: map[string]string{
			"userId":     user,
			"verifyCode": verifyCode,
		},
	}

	rsp, err := SendEmailToUser(req)
	if err != nil {
		return "", err
	}

	return rsp.DealId, nil

}
