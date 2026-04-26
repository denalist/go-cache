package main

import (
	"go-cache/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// receives from cache and return gin handler
func GetCacheValue(c *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Param("key")

		val, found := c.Get(key)

		if !found {
			ctx.JSON(http.StatusNotFound, gin.H{
				"key":   key,
				"found": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"key":   key,
			"value": val,
			"found": true,
		})
	}
}

// use json.Unmarshal parse Json encoded byte slide into Go Struct
type setRequest struct {
	Key        string      `json: "Key"`
	Value      interface{} `json: "value"`
	TTLSeconds int64       `json: "ttl_seconds"`
}

func SetCacheValue(c *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req setRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid JSON body",
			})
			return
		}

		ttl := time.Duration(req.TTLSeconds) * time.Second
		c.Set(req.Key, req.Value, ttl)
		ctx.JSON(http.StatusOK, gin.H{
			"key":     req.Key,
			"value":   req.Value,
			"Message": "Successfully POST.",
		})
	}
}

func DeleteCacheValue(c *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Param("key")
		res := c.Delete(key)
		if !res {
			ctx.JSON(http.StatusNotFound, gin.H{
				"key":     key,
				"deleted": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"key":     key,
			"deleted": true,
		})

	}
}
