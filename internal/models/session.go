package models

import "time"

type Session struct {
	UserId    int64
	SessionId string
	CreatedAt time.Time
}

func (s Session) GetUserId() int64 {
	return s.UserId
}

func (s Session) GetSessionId() string {
	return s.SessionId
}
