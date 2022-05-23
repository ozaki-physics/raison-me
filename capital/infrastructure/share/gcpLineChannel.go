package share

import "log"

func CreateCredentialLineGcp() CredentialLine {
	var ld lineDto

	secret, err := GcpSecretValue("LINE_CHANNEL_SECRET", 2)
	if err != nil {
		log.Fatalln(err)
	}
	ld.secret = secret

	token, err := GcpSecretValue("LINE_CHANNEL_TOKEN", 2)
	if err != nil {
		log.Fatalln(err)
	}
	ld.token = token

	return &ld
}
