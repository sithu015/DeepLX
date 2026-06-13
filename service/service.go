/*
 * @Author: Vincent Yang
 * @Date: 2023-07-01 21:45:34
 * @LastEditors: Jason Lyu
 * @LastEditTime: 2025-04-08 13:45:00
 * @FilePath: /DeepLX/main.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */

package service

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/OwO-Network/DeepLX/translate"
)

func authMiddleware(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Token != "" {
			providedToken := c.GetHeader("Authorization")

			// Compatibility with the Bearer token format
			if providedToken != "" {
				parts := strings.Split(providedToken, " ")
				if len(parts) == 2 {
					if parts[0] == "Bearer" || parts[0] == "DeepL-Auth-Key" {
						providedToken = parts[1]
					} else {
						providedToken = ""
					}
				} else {
					providedToken = ""
				}
			}

			// Fallback: check query parameters (Approach A)
			if providedToken == "" {
				providedToken = c.Query("token")
				if providedToken == "" {
					providedToken = c.Query("key")
				}
			}

			// Fallback: check path parameter (Approach B)
			if providedToken == "" {
				providedToken = c.Param("token")
			}

			// Security fix: use constant-time comparison to prevent timing side-channel attacks
			hProvided := sha256.Sum256([]byte(providedToken))
			hConfigured := sha256.Sum256([]byte(cfg.Token))
			tokenMatch := subtle.ConstantTimeCompare(hProvided[:], hConfigured[:]) == 1

			if !tokenMatch {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Invalid access token",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

type PayloadFree struct {
	TransText   string `json:"text"`
	SourceLang  string `json:"source_lang"`
	TargetLang  string `json:"target_lang"`
	TagHandling string `json:"tag_handling"`
}

type PayloadAPI struct {
	Text        []string `json:"text"`
	TargetLang  string   `json:"target_lang"`
	SourceLang  string   `json:"source_lang"`
	TagHandling string   `json:"tag_handling"`
}

func Router(cfg *Config) *gin.Engine {
	proxyURL := cfg.Proxy


	if cfg.Token != "" {
		fmt.Println("Access token is set.")
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Defining the root endpoint which returns the project details
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "DeepL Free API, Developed by sjlleo and missuo. Go to /translate with POST. http://github.com/OwO-Network/DeepLX",
		})
	})

	// Free translation handler logic
	freeHandler := func(c *gin.Context) {
		req := PayloadFree{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid request payload",
			})
			return
		}

		sourceLang := req.SourceLang
		targetLang := req.TargetLang
		translateText := req.TransText
		tagHandling := req.TagHandling

		if tagHandling != "" && tagHandling != "html" && tagHandling != "xml" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid tag_handling value. Allowed values are 'html' and 'xml'.",
			})
			return
		}

		result, err := translate.TranslateByDeepLX(c.Request.Context(), sourceLang, targetLang, translateText, tagHandling, proxyURL, "")
		if err != nil {
			log.Printf("Translation failed: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Translation failed",
			})
			return
		}

		if result.Code == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{
				"code":         http.StatusOK,
				"id":           result.ID,
				"data":         result.Data,
				"alternatives": result.Alternatives,
				"source_lang":  result.SourceLang,
				"target_lang":  result.TargetLang,
				"method":       result.Method,
			})
		} else {
			c.JSON(result.Code, gin.H{
				"code":    result.Code,
				"message": result.Message,
			})
		}
	}

	// Free API endpoints (both standard and path-based fallback for extensions)
	r.POST("/translate", authMiddleware(cfg), freeHandler)
	r.POST("/:token/translate", authMiddleware(cfg), freeHandler)

	// Pro translation handler logic
	proHandler := func(c *gin.Context) {
		req := PayloadFree{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid request payload",
			})
			return
		}

		sourceLang := req.SourceLang
		targetLang := req.TargetLang
		translateText := req.TransText
		tagHandling := req.TagHandling

		dlSession := cfg.DlSession

		if tagHandling != "" && tagHandling != "html" && tagHandling != "xml" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid tag_handling value. Allowed values are 'html' and 'xml'.",
			})
			return
		}

		cookie := c.GetHeader("Cookie")
		if cookie != "" {
			dlSession = strings.Replace(cookie, "dl_session=", "", -1)
		}

		if dlSession == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "No dl_session Found",
			})
			return
		}

		result, err := translate.TranslateByDeepLX(c.Request.Context(), sourceLang, targetLang, translateText, tagHandling, proxyURL, dlSession)
		if err != nil {
			log.Printf("Translation failed: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Translation failed",
			})
			return
		}

		if result.Code == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{
				"code":         http.StatusOK,
				"id":           result.ID,
				"data":         result.Data,
				"alternatives": result.Alternatives,
				"source_lang":  result.SourceLang,
				"target_lang":  result.TargetLang,
				"method":       result.Method,
			})
		} else {
			c.JSON(result.Code, gin.H{
				"code":    result.Code,
				"message": result.Message,
			})
		}
	}

	// Pro API endpoints (both standard and path-based fallback for extensions)
	r.POST("/v1/translate", authMiddleware(cfg), proHandler)
	r.POST("/v1/:token/translate", authMiddleware(cfg), proHandler)

	// Free API endpoint, Consistent with the official API format
	r.POST("/v2/translate", authMiddleware(cfg), func(c *gin.Context) {
		var translateText string
		var targetLang string

		translateText = c.PostForm("text")
		targetLang = c.PostForm("target_lang")

		if translateText == "" || targetLang == "" {
			var jsonData struct {
				Text       []string `json:"text"`
				TargetLang string   `json:"target_lang"`
			}

			if err := c.BindJSON(&jsonData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "Invalid request payload",
				})
				return
			}

			translateText = strings.Join(jsonData.Text, "\n")
			targetLang = jsonData.TargetLang
		}

		result, err := translate.TranslateByDeepLX(c.Request.Context(), "", targetLang, translateText, "", proxyURL, "")
		if err != nil {
			log.Printf("Translation failed: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Translation failed",
			})
			return
		}

		if result.Code == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{
				"translations": []map[string]interface{}{
					{
						"detected_source_language": result.SourceLang,
						"text":                     result.Data,
					},
				},
			})
		} else {
			c.JSON(result.Code, gin.H{
				"code":    result.Code,
				"message": result.Message,
			})
		}
	})

	// Catch-all route to handle undefined paths
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Path not found",
		})
	})

	return r
}
