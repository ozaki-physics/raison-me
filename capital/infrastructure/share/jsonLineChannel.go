package share

// dataLineChannel JSON から struct に変換する
type dataLineChannel struct {
	Data struct {
		Service struct {
			Secret string `json:"secret"`
			Token  string `json:"token"`
		} `json:"line_channel"`
	} `json:"data"`
}

func CreateCredentialLineJson() CredentialLine {
	var d dataLineChannel
	JsonToStruct(d, "./capital/infrastructure/share/json/key.json")
	return &lineDto{
		d.Data.Service.Secret,
		d.Data.Service.Token,
	}
}
