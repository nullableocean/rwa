package services

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"rwa/internal/models"
	"time"
)

const sessionIdLength = 16

type SessionRepository interface {
	Get(sessionId string) (*models.Session, error)
	GetAllByUser(userId int64) ([]*models.Session, error)
	Save(models.Session) error
	DeleteAllByUser(userId int64) error
	Delete(sessionId string) error
}

type SessionManager struct {
	sessionRepo SessionRepository
	userService *UserService
}

func NewSessionManager(sesRepo SessionRepository, us *UserService) *SessionManager {
	return &SessionManager{
		sessionRepo: sesRepo,
		userService: us,
	}
}

func (sm *SessionManager) Get(sessionId string) (*models.Session, error) {
	return sm.sessionRepo.Get(sessionId)
}

func (sm *SessionManager) GetAllByUser(user models.User) ([]*models.Session, error) {
	return sm.sessionRepo.GetAllByUser(user.ID)
}

func (sm *SessionManager) Create(user models.User) (*models.Session, error) {
	sessionId, err := sm.generateSessionId()
	if err != nil {
		return nil, err
	}

	session := models.Session{
		UserId:    user.ID,
		SessionId: sessionId,
		CreatedAt: time.Now(),
	}

	err = sm.sessionRepo.Save(session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (sm *SessionManager) Delete(sessionId string) error {
	return sm.sessionRepo.Delete(sessionId)
}

func (sm *SessionManager) DeleteAllByUser(user models.User) error {
	return sm.sessionRepo.DeleteAllByUser(user.ID)
}

func (sm *SessionManager) generateSessionId() (string, error) {
	b := make([]byte, sessionIdLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write(b)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
