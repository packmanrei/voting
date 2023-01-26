package controller

import (
	"crypto/sha256"
	"fmt"
	db "voting/models"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Account struct {
	db.Account
}

const (
	template = "views/templates/"
	index    = "views/index/"
)

func StartServer() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles("views/**/*.html")

	// import css file
	router.StaticFile("/style.css", "views/css/style.css")
	// import image file
	router.StaticFile("/favicon.png", "images/favicon.png")
	router.StaticFile("/login.png", "images/login.png")
	router.StaticFile("/register.png", "images/register.png")
	router.StaticFile("/logout.png", "images/logout.png")
	router.StaticFile("/account.png", "images/account.png")
	router.StaticFile("/contact.png", "images/contact.png")
	// import js file
	router.StaticFile("/main.js", "views/js/main.js")

	// create render
	router.HTMLRender = createMyRender()

	// session settings
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	homeRoute(router, "/")
	errorRoute(router, "/error")
	registerRoute(router, "/register")
	loginRoute(router, "/login")
	logoutRoute(router, "/logout")
	votingRoute(router, "/voting")
	postRoute(router, "/post")
	accountSettingRoute(router, "/account-settings")
	contactRoute(router, "/contact")

	return router
}

func createMyRender() multitemplate.Renderer {
	render := multitemplate.NewRenderer()

	// account
	render.AddFromFiles("login", template+"account.html", index+"login.html")
	render.AddFromFiles("register", template+"account.html", index+"register.html")
	render.AddFromFiles("createAccount", template+"account.html", index+"create_account.html")
	render.AddFromFiles("accountSettings", template+"account.html", index+"account_settings.html")
	render.AddFromFiles("updateSettings", template+"account.html", index+"update_settings.html")
	render.AddFromFiles("deleteAccount", template+"account.html", index+"delete_account.html")
	render.AddFromFiles("goodbye", template+"account.html", index+"goodbye.html")
	render.AddFromFiles("contact", template+"account.html", index+"contact.html")
	render.AddFromFiles("thanks", template+"account.html", index+"thanks.html")
	render.AddFromFiles("error", template+"account.html", index+"error.html")
	// home
	render.AddFromFiles("homeNotLogedIn", template+"home.html", index+"header_not_logedin.html", index+"votings.html", index+"rooms.html")
	render.AddFromFiles("homeLogedIn", template+"home.html", index+"header_logedin.html", index+"votings.html", index+"rooms.html")
	render.AddFromFiles("post", template+"home.html", index+"header_logedin.html", index+"post.html", index+"rooms.html")
	render.AddFromFiles("voteLogedIn", template+"home.html", index+"header_logedin.html", index+"vote.html", index+"rooms.html")
	render.AddFromFiles("voteNotLogedIn", template+"home.html", index+"header_not_logedin.html", index+"vote.html", index+"rooms.html")
	render.AddFromFiles("voteHasVoted", template+"home.html", index+"header_logedin.html", index+"has_voted.html", index+"rooms.html")
	render.AddFromFiles("votePoster", template+"home.html", index+"header_logedin.html", index+"vote_poster.html", index+"rooms.html")

	return render
}

func toEncrypt(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}
