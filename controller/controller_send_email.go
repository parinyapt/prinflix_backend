package controller

import (
	"os"
	"strconv"
	"strings"

	modelController "github.com/parinyapt/prinflix_backend/model/controller"
	modelDatabase "github.com/parinyapt/prinflix_backend/model/database"
	utilsConfigFile "github.com/parinyapt/prinflix_backend/utils/config_file"
	"github.com/pkg/errors"

	PTGUmail "github.com/parinyapt/golang_utils/mail/v1"
)

func SendEmail(param modelController.ParamSendEmail) (err error) {
	var url string
	var subject string
	var body string

	switch param.Type {
	case modelDatabase.TemporaryCodeTypePasswordReset:
		url = utilsConfigFile.GetFrontendBaseURL() + utilsConfigFile.GetRedirectPagePath(utilsConfigFile.ResetPasswordPagePath)
		url = strings.Replace(url, ":session_id", param.Data, -1)
		subject = "Prinflix Reset Password"
		body = "Reset your password by clicking this link: " + url + ". This link will be expired in 15 minutes."
	default:
		return errors.New("[Controller][SendEmail()]->invalid type")
	}

	hostport, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		return errors.Wrap(err, "[Controller][SendEmail()]->fail to convert port to int")
	}
	err = PTGUmail.SendMail(PTGUmail.ParamConfigSendMail{
		SMTP: PTGUmail.ParamConfigSendMailSMTP{
			Host:     os.Getenv("EMAIL_SMTP_HOST"),
			Port:     hostport,
			Username: os.Getenv("EMAIL_SMTP_USERNAME"),
			Password: os.Getenv("EMAIL_SMTP_PASSWORD"),
		},
		From: PTGUmail.ParamConfigSendMailFrom{
			AliasName: os.Getenv("EMAIL_SMTP_FROM_ALIAS_NAME"),
			Email:     os.Getenv("EMAIL_SMTP_FROM_EMAIL"),
		},
		To: PTGUmail.ParamConfigSendMailTo{
			Email:    []string{param.Email},
			Subject:  subject,
			BodyType: PTGUmail.BodyTypePlain,
			Body:     body,
		},
	})
	if err != nil {
		return errors.Wrap(err, "[Controller][SendEmail()]->fail to send email")
	}

	return nil
}