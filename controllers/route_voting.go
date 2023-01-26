package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	db "voting/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func homeRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		//check session
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)
		votings := db.GetAllVotings()
		if check {
			account := db.SearchAccountByEmail(email)
			ctx.HTML(http.StatusOK, "homeLogedIn", gin.H{
				"account": account,
				"votings": votings,
			})
		} else {
			ctx.HTML(http.StatusOK, "homeNotLogedIn", gin.H{
				"votings": votings,
			})
		}
	})

	router.GET(path+"/:roomName", func(ctx *gin.Context) {
		templateName := "homeNotLogedIn"
		roomName := ctx.Param("roomName")
		votings := db.GetSingleRoomVotings(roomName)
		session := sessions.Default(ctx)
		_, check := session.Get("account_email").(string)
		if check {
			templateName = "homeLogedIn"
		}
		ctx.HTML(http.StatusOK, templateName, gin.H{
			"votings": votings,
		})
	})
}

func votingRoute(router *gin.Engine, path string) {

	router.GET(path+"/:id", func(ctx *gin.Context) {
		templateName := "voteNotLogedIn"

		// get voting
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.Redirect(http.StatusFound, "/error")
			return
		}

		// check session
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)
		account := db.SearchAccountByEmail(email)
		if check {
			templateName = "voteLogedIn"
		}
		voting := db.GetSingleVotingByID(id)
		if voting.ID == 0 {
			ctx.Redirect(http.StatusFound, "/error")
			return
		}
		posterAccount := db.SearchAccountByID(voting.PosterID)

		// has voted check
		voters := strings.Split(voting.VotersID, ",")
		var hasVoted bool = false
		for _, voterID := range voters {
			if id, err := strconv.Atoi(voterID); id == int(account.ID) {
				if err != nil {
					break
				}
				hasVoted = true
				templateName = "voteHasVoted"
				break
			}
		}

		// has watched check
		watchers := strings.Split(voting.Watchers, ",")
		var hasWatched bool = false
		for _, watcherID := range watchers {
			if id, err := strconv.Atoi(watcherID); id == int(account.ID) {
				if err != nil {
					break
				}
				hasWatched = true
				templateName = "voteHasVoted"
				break
			}
		}

		// check poster
		var isPoster bool = false
		if strconv.Itoa(int(account.ID)) == strconv.Itoa(voting.PosterID) {
			isPoster = true
			templateName = "votePoster"
		}

		// check condition
		var matchCond bool = false
		var sexCheck bool = false
		var ageCheck bool = false
		var textSex string
		var textAge string
		condition := db.GetCondition(id)
		// check sex
		if condition.Sex == "both" {
			sexCheck = true
			textSex = "指定なし"
		} else {
			if account.Sex == condition.Sex {
				sexCheck = true
			}
			if condition.Sex == "man" {
				textSex = "男性のみ"
			} else {
				textSex = "女性のみ"
			}
		}
		//check age
		switch condition.Age {
		case "more":
			if account.Age > condition.AgeNum {
				ageCheck = true
			}
			textAge = fmt.Sprintf("%d歳以上", condition.AgeNum)
		case "less":
			if account.Age < condition.AgeNum {
				ageCheck = true
			}
			textAge = fmt.Sprintf("%d歳以下", condition.AgeNum)
		default:
			ageCheck = true
			textAge = "指定なし"
		}
		if check {
			if sexCheck && ageCheck {
				matchCond = true
			} else if !isPoster {
				hasWatched = true
				templateName = "voteHasVoted"
			}
		}

		// route
		if hasVoted || isPoster || hasWatched {
			choices := strings.Split(voting.Choices, ",")
			votes := strings.Split(voting.Votes, ",")
			var sum int
			// results
			for _, vote := range votes {
				num, _ := strconv.Atoi(vote)
				sum += num
			}
			f_sum := float64(sum)
			var votingResults []db.VotingResult
			for i, vote := range votes {
				var votingResult db.VotingResult
				var result float64
				// calculate result
				if f_sum != 0 {
					num, _ := strconv.Atoi(vote)
					f_num := float64(num)
					result = math.Round(f_num * 100 / f_sum)
				} else {
					result = 0
				}
				votingResult.Result = int(result)
				votingResult.Choice = choices[i]
				votingResult.Vote = votes[i]
				votingResults = append(votingResults, votingResult)
			}

			if isPoster {
				url := fmt.Sprintf("%s/%d/delete", path, id)
				ctx.HTML(http.StatusOK, templateName, gin.H{
					"voting":        voting,
					"posterAccount": posterAccount,
					"votingResults": votingResults,
					"url":           url,
					"total":         f_sum,
					"textSex":       textSex,
					"textAge":       textAge,
				})
			} else {
				ctx.HTML(http.StatusOK, templateName, gin.H{
					"voting":        voting,
					"posterAccount": posterAccount,
					"votingResults": votingResults,
					"total":         f_sum,
					"textSex":       textSex,
					"textAge":       textAge,
				})
			}
		} else if matchCond || (!isPoster && !matchCond && !check) {
			choices := strings.Split(voting.Choices, ",")
			url := fmt.Sprintf("%s/%d/just-watch", path, id)
			ctx.HTML(http.StatusOK, templateName, gin.H{
				"voting":        voting,
				"posterAccount": posterAccount,
				"choices":       choices,
				"url":           url,
				"textSex":       textSex,
				"textAge":       textAge,
			})
		}
	})

	router.POST(path+"/:id", func(ctx *gin.Context) {
		number, _ := strconv.Atoi(ctx.PostForm("choice"))
		id, _ := strconv.Atoi(ctx.Param("id"))
		voting := db.GetSingleVotingByID(id)
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)

		// session check
		if !check {
			session.Set("needLoginMessage", "投票にはログインが必要です")
			session.Save()
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
		// write votersID
		account := db.SearchAccountByEmail(email)
		voters := strings.Split(voting.VotersID, ",")
		if voters[0] == "" {
			voters[0] = strconv.Itoa(int(account.ID))
		} else {
			voters = append(voters, strconv.Itoa(int(account.ID)))
		}
		var votersID string
		for i, voterID := range voters {
			if i != 0 {
				votersID += ","
			}
			votersID += voterID
		}
		voting.VotersID = votersID

		// vote + 1
		votesSlice := strings.Split(voting.Votes, ",")
		singleVote, _ := strconv.Atoi(votesSlice[number-1])
		singleVote += 1
		votesSlice[number-1] = strconv.Itoa(singleVote)
		var votes string
		for i, vote := range votesSlice {
			if i != 0 {
				votes += ","
			}
			votes += vote
		}
		voting.Votes = votes
		db.UpdateSingleVoting(voting)

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/voting/"+strconv.Itoa(id)))
	})

	// Delete Post
	router.GET(path+"/:id/delete", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.Redirect(http.StatusFound, "/")
			return
		}
		voting := db.GetSingleVotingByID(id)
		db.DeleteVoting(voting)
		ctx.Redirect(http.StatusFound, "/")
	})

	// Just Watch
	router.GET(path+"/:id/just-watch", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		voting := db.GetSingleVotingByID(id)
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)

		// session check
		if !check {
			session.Set("needLoginMessage", "結果を見るにはログインが必要です")
			session.Save()
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
		// write votersID
		account := db.SearchAccountByEmail(email)
		watchers := strings.Split(voting.Watchers, ",")
		if watchers[0] == "" {
			watchers[0] = strconv.Itoa(int(account.ID))
		} else {
			watchers = append(watchers, strconv.Itoa(int(account.ID)))
		}
		var watchersID string
		for i, watcherID := range watchers {
			if i != 0 {
				watchersID += ","
			}
			watchersID += watcherID
		}
		voting.Watchers = watchersID

		db.UpdateSingleVoting(voting)

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/voting/"+strconv.Itoa(id)))
	})
}

func postRoute(router *gin.Engine, path string) {

	router.GET(path, func(ctx *gin.Context) {
		//check session
		session := sessions.Default(ctx)
		_, check := session.Get("account_email").(string)
		if !check {
			session.Set("needLoginMessage", "投稿にはログインが必要です")
			session.Save()
			ctx.Redirect(http.StatusFound, "/login")
			return
		}
		ctx.HTML(http.StatusOK, "post", nil)
	})

	router.POST(path, func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		email, check := session.Get("account_email").(string)
		if !check {
			ctx.Redirect(http.StatusFound, "/error")
			return
		}
		// slice -> string
		choicesSlice := ctx.PostFormArray("choice")
		nilCounter := 0
		for _, choice := range choicesSlice {
			if choice == "" {
				nilCounter++
			}
		}
		choicesSlice = choicesSlice[:len(choicesSlice)-nilCounter]
		var choices string
		var votes string
		for i, choice := range choicesSlice {
			if i != 0 {
				choices = choices + ","
				votes = votes + ","
			}
			choices = choices + choice
			votes = votes + "0"
		}

		account := db.SearchAccountByEmail(email)
		accountID := int(account.ID)
		roomName := ctx.PostForm("room")
		title := ctx.PostForm("title")
		content := ctx.PostForm("content")
		// create post
		voting := db.CreatePost(accountID, roomName, title, content, choices, votes)

		// condition
		sex := ctx.PostForm("sex")
		age := ctx.PostForm("age")
		ageNum, err := strconv.Atoi(ctx.PostForm("age_num"))
		if err != nil {
			ctx.Redirect(http.StatusFound, "/error")
		}
		db.AddCondition(int(voting.ID), sex, age, ageNum)

		ctx.Redirect(http.StatusFound, "/")
	})
}
