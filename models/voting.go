package db

func CreatePost(posterID int, roomName string, title string, content string, choices string, votes string) Voting {
	db := OpenDB()
	voting := Voting{PosterID: posterID, RoomName: roomName, Title: title, Content: content, Choices: choices, Votes: votes}
	db.Create(&voting)
	return voting
}

func GetAllVotings() []Voting {
	db := OpenDB()
	var votings []Voting
	db.Order("id desc").Find(&votings)
	return votings
}

func GetSingleRoomVotings(RoomName string) []Voting {
	db := OpenDB()
	var votings []Voting
	db.Where("room_name = ?", RoomName).Order("id desc").Find(&votings)
	return votings
}

func GetSingleVotingByID(id int) Voting {
	db := OpenDB()
	var voting Voting
	db.Where("id = ?", id).First(&voting)
	return voting
}

func UpdateSingleVoting(voting Voting) {
	db := OpenDB()
	db.Model(&Voting{}).Where("id = ?", voting.ID).Updates(voting)
}

func DeleteVoting(voting Voting) {
	db := OpenDB()
	db.Delete(&voting)
}

func AddCondition(votingID int, sex string, age string, ageNum int) {
	db := OpenDB()
	condition := Condition{VotingID: votingID, Sex: sex, Age: age, AgeNum: ageNum}
	db.Create(&condition)
}

func GetCondition(votingID int) Condition {
	db := OpenDB()
	var condition Condition
	db.Where("voting_id = ?", votingID).First(&condition)
	return condition
}
