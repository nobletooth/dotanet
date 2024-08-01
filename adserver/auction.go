package main

import (
	"common"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetImage(adID uint) (string, error) {
	url := fmt.Sprintf("%v/ads/%d/picture", *PanelUrlPicUrl, adID)
	return url, nil
}

func GetAdHandler(c *gin.Context) {
	pubID := c.Param("pubID")
	if len(allAds) == 0 {
		sendAdResponse(c, nil, pubID)
		return
	}

	var newAds []common.AdWithMetrics
	var experiencedAds []common.AdWithMetrics
	var ctrPrices []float64

	for _, ad := range allAds {
		if ad.ImpressionCount < *NewAdImpressionThreshold {
			newAds = append(newAds, ad)
		} else {
			ctrPrice := float64(ad.ClickCount) / float64(ad.ImpressionCount) * ad.Price
			ctrPrices = append(ctrPrices, ctrPrice)
			experiencedAds = append(experiencedAds, ad)
		}
	}

	rand.Seed(time.Now().UnixNano())
	selectNewAd := rand.Float64() < *NewAdSelectionProbability

	var finalAd *common.AdWithMetrics
	if selectNewAd && len(newAds) > 0 {
		rand.Seed(time.Now().UnixNano())
		selectedAd := newAds[rand.Intn(len(newAds))]
		finalAd = &selectedAd
	} else if len(experiencedAds) > 0 {
		totalScore := 0.0
		for _, ctrPrice := range ctrPrices {
			totalScore += ctrPrice
		}

		rand.Seed(time.Now().UnixNano())
		randomPoint := rand.Float64() * totalScore
		currentSum := 0.0
		for i, ad := range experiencedAds {
			currentSum += ctrPrices[i]
			if randomPoint <= currentSum {
				finalAd = &ad
				break
			}
		}
	} else {
		finalAd = nil
	}
	sendAdResponse(c, finalAd, pubID)
}

func sendAdResponse(c *gin.Context, ad *common.AdWithMetrics, pubID string) {
	if ad == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No ads available"})
		return
	}
	fmt.Printf("adID: %d, ad title: %s, ad price: %f\n", ad.Id, ad.Title, ad.Price)
	imageDataurl, err := GetImage(ad.AdInfo.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	publisherID, err := strconv.ParseUint(pubID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	var impression = common.ViewedEvent{
		ID:   uuid.New(),
		Pid:  int(publisherID),
		AdId: int(ad.Id),
	}

	var click = common.ClickedEvent{
		ID:           uuid.New(),
		Pid:          int(publisherID),
		AdId:         int(ad.Id),
		ImpressionID: impression.ID,
	}

	// encryption
	jsonClickParams, _ := json.Marshal(click)
	encryptedClickParams, err := encrypt(jsonClickParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt click params"})
		return
	}

	jsonImpressionParams, _ := json.Marshal(impression)
	encryptedImpressionParams, err := encrypt(jsonImpressionParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt impression params"})
		return
	}

	response := gin.H{
		"Title":          ad.Title,
		"ImageData":      imageDataurl,
		"ClicksURL":      fmt.Sprintf("%v/click/%v", *EventServiceUrl, encryptedClickParams),
		"ImpressionsURL": fmt.Sprintf("%v/impression/%v", *EventServiceUrl, encryptedImpressionParams),
	}
	c.JSON(http.StatusOK, response)
}

func encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher([]byte(*secretKey))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
