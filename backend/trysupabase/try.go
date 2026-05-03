package trysupabase

func Try() string {
	return transactionPooler_ReadUsers_globalConfig()
}

func TryAll() {
	getRows()
	getRowsV2()
	getColumnsWhere(0)
	insertRow("こんにちは")
	getRowsSchema()
	insertAppUsersRow("テストユーザー02")

	transactionPooler_Connect()
	sessionPooler_Connect()
	transactionPooler_ReadUsers()
	sessionPooler_ReadUsers()
	transactionPooler_ReadUsers_keyJSON()
	transactionPooler_ReadUsers_globalConfig()
}
