package settings

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Bryan-an/tasker-backend/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

type replaceInput struct {
	Notifications *notification `json:"notifications" binding:"required"`
	Theme         *string       `json:"theme" binding:"required,oneof=dark light"`
}

func (h handler) ReplaceSettings(c *gin.Context) {
	uid, err := utils.ExtractTokenID(c)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var input replaceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			out := utils.FillErrors(ve)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		} else {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		return
	}

	settingsCollection := h.DB.Collection("settings")
	filter := bson.D{{Key: "user_id", Value: uid}}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "notifications", Value: input.Notifications},
				{Key: "theme", Value: input.Theme},
				{Key: "updated_at", Value: time.Now()},
			},
		},
	}

	result, err := settingsCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf(
				"settings not found for user with id '%s'",
				uid.Hex(),
			),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "settings replaced successfully",
	})
}
