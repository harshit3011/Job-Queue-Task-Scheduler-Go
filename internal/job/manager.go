package job

import (
	"time"

	"github.com/gin-gonic/gin"
)


func ManageJos(job Job) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}

func IsAllowedTime() bool{
	now := time.Now()

	allowedStart:= time.Date(now.Year(),now.Month(),now.Day(),9,0,0,0,now.Location())
	allowedEnd:=time.Date(now.Year(),now.Month(),now.Day(),18,0,0,0,now.Location())

	if !now.Before(allowedStart) && !now.After(allowedEnd){
		return true
	} else{
		return false
	}
}