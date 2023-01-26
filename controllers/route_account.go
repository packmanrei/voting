package controller

import (
	"fmt"
	"net/http"
	"strconv"

	db "voting/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func registerRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register", nil)
	})

	router.POST(path, func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		used := db.CheckAccount(email)
		session := sessions.Default(ctx)
		if used {
			ctx.Redirect(http.StatusFound, path)
			session.Set("errMessage", "このEメールは既に使用されています")
			session.Save()
			return
		}
		session.Set("account_email", email)
		session.Save()
		db.CreateAccount(email)
		account := db.SearchAccountByEmail(email)
		id := strconv.Itoa(int(account.ID))
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/register/"+id))
	})

	router.GET(path+"/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.HTML(http.StatusOK, "createAccount", gin.H{
			"url": fmt.Sprintf(path + "/" + id),
		})
	})

	router.POST(path+"/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		username := ctx.PostForm("username")
		sex := ctx.PostForm("sex")
		age, err := strconv.Atoi(ctx.PostForm("age"))
		if err != nil {
			ctx.Redirect(http.StatusFound, fmt.Sprintf("/register/"+id))
			return
		}

		password := ctx.PostForm("password")
		passwordAgein := ctx.PostForm("password_again")
		if password != passwordAgein {
			ctx.Redirect(http.StatusFound, fmt.Sprintf("/register/"+id))
			return
		}
		password = toEncrypt(password)
		session := sessions.Default(ctx)
		email := session.Get("account_email").(string)
		account := db.Account{Username: username, Password: password, Sex: sex, Age: age, Email: email}
		db.FirstSettings(account)
		ctx.Redirect(http.StatusFound, "/")
	})
}

func loginRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		errMessage, _ := session.Get("errMessage").(string)
		needLoginMessage, _ := session.Get("needLoginMessage").(string)
		ctx.HTML(http.StatusOK, "login", gin.H{
			"errMessage":       errMessage,
			"needLoginMessage": needLoginMessage,
		})
	})

	router.POST(path, func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Clear()
		session.Save()
		email := ctx.PostForm("email")
		password := toEncrypt(ctx.PostForm("password"))
		account := db.SearchAccount(password, email)
		if account.ID == 0 {
			ctx.Redirect(http.StatusFound, path)
			session.Set("errMessage", "Eメールまたはパスワードが間違っています")
			session.Save()
		} else {
			session := sessions.Default(ctx)
			session.Clear()
			session.Set("account_email", account.Email)
			session.Save()
			ctx.Redirect(http.StatusFound, "/")
		}
	})
}

func logoutRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Clear()
		err := session.Save()
		db.CheckError(err)
		ctx.Redirect(http.StatusFound, "/")
	})
}

func accountSettingRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)
		if !check {
			ctx.Redirect(http.StatusFound, "/error")
			return
		}
		account := db.SearchAccountByEmail(email)
		if account.Sex == "man" {
			account.Sex = "男性"
		} else {
			account.Sex = "女性"
		}
		ctx.HTML(http.StatusOK, "accountSettings", gin.H{
			"account": account,
		})
	})

	router.GET(path+"/delete-account", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "deleteAccount", gin.H{})
	})

	router.GET("/goodbye", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)
		if !check {
			ctx.Redirect(http.StatusFound, "/error")
			return
		}
		account := db.SearchAccountByEmail(email)
		db.DeleteAccount(account)
		session.Clear()
		session.Save()
		ctx.HTML(http.StatusOK, "goodbye", nil)
	})

	router.GET(path+"/:kind", func(ctx *gin.Context) {
		kind := ctx.Param("kind")
		url := path + "/" + kind
		session := sessions.Default(ctx)
		_, check := session.Get("account_email").(string)
		if !check {
			ctx.Redirect(http.StatusFound, "/error")
			return
		}
		errMessage, _ := session.Get("errMessage").(string)
		session.Delete("errMessage")
		session.Save()
		var message string
		switch kind {
		case "username":
			message = "新しいユーザー名を入力してください"
		case "email":
			message = "新しいEメールを入力してください"
		case "password":
			message = "新しいパスワードを入力してください"
		default:
			ctx.Redirect(http.StatusFound, "/account-settings")
			return
		}
		ctx.HTML(http.StatusOK, "updateSettings", gin.H{
			"url":        url,
			"message":    message,
			"errMessage": errMessage,
		})
	})

	router.POST(path+"/:kind", func(ctx *gin.Context) {
		kind := ctx.Param("kind")
		input := ctx.PostForm("input")
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)
		if !check {
			ctx.Redirect(http.StatusFound, "/error")
			return
		}
		account := db.SearchAccountByEmail(email)
		switch kind {
		case "username":
			account.Username = input
		case "email":
			used := db.CheckAccount(input)
			if used {
				ctx.Redirect(http.StatusFound, "/account-settings/email")
				session.Set("errMessage", "入力されたEメールは既に使用されています")
				session.Save()
				return
			}
			account.Email = input
			session.Set("account_email", input)
			session.Save()
		case "password":
			pass := toEncrypt(input)
			account.Password = pass
		}
		db.UpdateAccount(account)
		ctx.Redirect(http.StatusFound, "/account-settings")
	})
}

func contactRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "contact", nil)
	})

	router.POST(path, func(ctx *gin.Context) {
		content := ctx.PostForm("content")
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)
		if !check {
			ctx.Redirect(http.StatusFound, "error")
			return
		}
		db.ContactToUs(email, content)
		ctx.Redirect(http.StatusFound, "/thanks")
	})

	router.GET("/thanks", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "thanks", nil)
	})
}

func errorRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Clear()
		session.Save()
		ctx.HTML(http.StatusOK, "error", nil)
	})
}
