// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/core/controller"
	"github.com/axetroy/go-server/core/exception"
	"github.com/axetroy/go-server/core/helper"
	"github.com/axetroy/go-server/core/middleware"
	"github.com/axetroy/go-server/core/model"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

// Get notification detail
func Get(c controller.Context, id string) (res schema.Response) {
	var (
		err  error
		data schema.Notification
		tx   *gorm.DB
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, err)
	}()

	tx = database.Db.Begin()

	notificationInfo := model.Notification{}
	NotificationMark := model.NotificationMark{Id: id}

	sql := fmt.Sprintf("LEFT JOIN notification_mark ON notification_mark.id = notification.id AND notification.id = '%s'", notificationInfo.Id)

	if err = tx.Table(notificationInfo.TableName()).Joins(sql).Last(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if er := mapstructure.Decode(notificationInfo, &data.NotificationPure); er != nil {
		err = er
		return
	}

	if err = tx.Where(&NotificationMark).Last(&NotificationMark).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			data.NotificationPure.Read = false
			err = nil
		} else {
			return
		}
	} else {
		data.NotificationPure.Read = NotificationMark.Read
		data.ReadAt = NotificationMark.CreatedAt.Format(time.RFC3339Nano)
	}

	data.CreatedAt = notificationInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = notificationInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

// GetRouter get notification detail router
func GetRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	id := c.Param("id")

	res = Get(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, id)
}
