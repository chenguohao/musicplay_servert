package dao

type PlayCountModel struct {
	Model
	PlaySender string `json:"play_sender"`
	PlayTarget int64  `json:"play_target"`
}
