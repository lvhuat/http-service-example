package msgapi

import (
	"github.com/lworkltd/kits/service/invoke"
	invokeutils "github.com/lworkltd/kits/utils/invoke"
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

func SendEmailToUser(req *SendMailRequest) error {
	httpRsp, err := invoke.Addr(seviceAddr).
		Post("/msgapi/v1/mail/send").
		Hystrix(100000, 0, 0).
		Json(req).
		Response()

	var dataRsp SendMailResponse
	cerr := invokeutils.ExtractHttpResponse(seviceAddr, err, httpRsp, &dataRsp)
	if cerr != nil {
		return cerr
	}

	return nil

}

func SendRegisterVerifyMail(user, email, verifyCode string) error {
	req := &SendMailRequest{
		Email:      email,
		TemplateId: "REGISTER_VERIFY_CODE",
		Parameters: map[string]string{
			"userId":     user,
			"verifyCode": verifyCode,
		},
	}

	return SendEmailToUser(req)
}
