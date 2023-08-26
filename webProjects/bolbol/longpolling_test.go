package bolbol

import (
	"context"
	"github.com/labstack/echo/v4"
	"own/example/bolbol/entity"
	"strconv"
	"time"
)

/*
	https://medium.com/@mhrlife/building-a-real-time-notification-service-with-golang-golang-basics-92e0dcb48d4d
*/

func (s *Server) listen(c echo.Context) error {
	clientId, _ := strconv.Atoi(c.Param("id"))
	notifications, err := s.bolbol.GetNotifications(c.Request().Context(), clientId)
	if err != nil {
		return err
	}
	return c.JSON(200, notifications)
}

type NotifyRequest struct {
	UserID int `json:"userID"`

	UnreadMessage     *entity.UnreadMessagesNotification `json:"unreadMessage"`
	UnreadWorkRequest *entity.UnreadWorkRequest          `json:"unreadWorkRequest"`
}

func (n *NotifyRequest) Notification() entity.Notification {
	if n.UnreadMessage != nil {
		return n.UnreadMessage
	}
	if n.UnreadWorkRequest != nil {
		return n.UnreadWorkRequest
	}
	panic("bad notification")
}

func (s *Server) notify(c echo.Context) error {
	var request NotifyRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	// the notify method's timeout doesn't depend on the request's timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.bolbol.Notify(ctx, request.UserID, request.Notification()); err != nil {
		return err
	}
	return c.String(201, "notification created")
}
