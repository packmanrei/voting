package db

func CreateAccount(email string) {
	db := OpenDB()
	account := Account{Email: email}
	db.Create(&account)
}

func SearchAccountByEmail(email string) Account {
	db := OpenDB()
	var account Account
	db.Where("email = ?", email).First(&account)
	return account
}

func SearchAccountByID(id int) Account {
	db := OpenDB()
	var account Account
	db.Where("id = ?", id).First(&account)
	return account
}

func CheckAccount(email string) bool {
	db := OpenDB()
	account := Account{Email: email}
	db.Where("email = ?", email).First(&account)
	if account.ID != 0 {
		return true
	} else {
		return false
	}
}

func FirstSettings(account Account) {
	db := OpenDB()
	db.Model(&Account{}).Where("email = ?", account.Email).Updates(account)
}

func SearchAccount(password string, email string) Account {
	db := OpenDB()
	var account Account
	db.Where("password = ? AND email = ?", password, email).First(&account)
	return account
}

func UpdateAccount(account Account) {
	db := OpenDB()
	db.Model(&Account{}).Where("id = ?", account.ID).Updates(account)
}

func DeleteAccount(account Account) {
	db := OpenDB()
	db.Delete(&account)
}

func ContactToUs(email string, content string) {
	db := OpenDB()
	db.Create(&Contact{Email: email, Content: content})
}
