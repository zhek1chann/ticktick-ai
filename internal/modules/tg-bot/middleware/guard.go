package middleware

import tele "gopkg.in/telebot.v3"

type Guard struct {
	adminIDs map[int64]struct{}
}

func New(adminIDs []int64) *Guard {
	ids := make(map[int64]struct{}, len(adminIDs))
	for _, id := range adminIDs {
		ids[id] = struct{}{}
	}
	return &Guard{adminIDs: ids}
}

func (g *Guard) RequireAdmin(c tele.Context, next tele.HandlerFunc) error {
	sender := c.Sender()
	if sender == nil {
		return c.Send("Доступ запрещен.")
	}
	if _, ok := g.adminIDs[sender.ID]; !ok {
		return c.Send("Доступ запрещен.")
	}
	return next(c)
}
