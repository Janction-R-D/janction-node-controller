package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxWalletAddressKey = "wallet_address"

var ErrorUserNotLogin = errors.New("user not login")

func getCurrentWalletAddress(c *gin.Context) (address string, err error) {
	address_, ok := c.Get(CtxWalletAddressKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	address = address_.(string)
	return
}

func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
