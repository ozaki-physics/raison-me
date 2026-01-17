package trysupabase

func Try() {
	transactionPooler_ReadUsers()
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
