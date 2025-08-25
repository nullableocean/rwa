package ram

import (
	"rwa/internal/models"
	"sync"
)

type SessionRepository struct {
	store           map[string]models.Session
	userSessionsMap map[int64][]string
	mu              *sync.Mutex
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		store:           make(map[string]models.Session),
		userSessionsMap: make(map[int64][]string),
		mu:              &sync.Mutex{},
	}
}

func (r *SessionRepository) Get(sessionId string) (*models.Session, error) {
	s, exist := r.store[sessionId]
	if !exist {
		return nil, models.ErrNotFound
	}

	return &s, nil
}
func (r *SessionRepository) GetAllByUser(userId int64) ([]*models.Session, error) {
	m, exist := r.userSessionsMap[userId]

	if !exist || len(m) == 0 {
		return []*models.Session{}, models.ErrNotFound
	}

	sessions := make([]*models.Session, 0, len(m))
	for _, sesId := range m {
		ses := r.store[sesId]
		sessions = append(sessions, &ses)
	}

	return sessions, nil
}
func (r *SessionRepository) Save(session models.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[session.SessionId] = session

	if _, ex := r.userSessionsMap[session.UserId]; !ex {
		r.userSessionsMap[session.UserId] = make([]string, 0)
	}

	r.userSessionsMap[session.UserId] = append(r.userSessionsMap[session.UserId], session.SessionId)

	return nil
}
func (r *SessionRepository) DeleteAllByUser(userId int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, sId := range r.userSessionsMap[userId] {
		delete(r.store, sId)
	}

	delete(r.userSessionsMap, userId)

	return nil
}
func (r *SessionRepository) Delete(sessionId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exist := r.store[sessionId]
	if !exist {
		return models.ErrNotFound
	}

	delete(r.store, sessionId)

	uId := session.UserId
	var ind int
	for i, sId := range r.userSessionsMap[uId] {
		if sId == sessionId {
			ind = i
			break
		}
	}

	r.userSessionsMap[uId] = append(r.userSessionsMap[uId][ind:], r.userSessionsMap[uId][:ind+1]...)

	return nil
}
