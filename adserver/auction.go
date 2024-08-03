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
	"slices"

	_ "slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


var RandomGenerator randomGenerator = &defaultRandomGenerator{}

type randomGenerator interface {
	Float64() float64
	Intn(n int) int
}

type defaultRandomGenerator struct{}

func (r *defaultRandomGenerator) Float64() float64 {
	return rand.Float64()
}

func (r *defaultRandomGenerator) Intn(n int) int {
	return rand.Intn(n)
}

func GetImage(adID uint) (string, error) {
	url := fmt.Sprintf("%v/ads/%d/picture", *PanelUrlPicUrl, adID)
	return url, nil
}

func GetAdHandler(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	pubID := c.Param("pubID")

	// set cookie
	_, err := c.Cookie("userId")
	if err != nil {
		userID := uuid.New().String()
		c.SetCookie("userId", userID, 24*180*60*60, "/", "", false, true)
	}

	if len(allAds) == 0 {
		sendAdResponse(c, nil, pubID)
		return
	}
	pubidint, _ := strconv.Atoi(pubID)

	var newAds []common.AdWithMetrics
	var experiencedAds []common.AdWithMetrics
	var ctrPrices []float64
	var acceptbleadd []common.AdWithMetrics
	for _, firstad := range allAds {
		if len(firstad.PreferdPubID) == 0 {
			continue
		} else if slices.Contains(firstad.PreferdPubID, uint(pubidint)) || firstad.PreferdPubID[0] == 0 {
			acceptbleadd = append(acceptbleadd, firstad)
		}

	}
	fmt.Println(acceptbleadd)

	for _, ad := range acceptbleadd {
		if ad.ImpressionCount < *NewAdImpressionThreshold {
			newAds = append(newAds, ad)
		} else {
			if ad.Price > 0 {
				ctrPrice := float64(ad.ClickCount) / float64(ad.ImpressionCount) * ad.Price
				ctrPrices = append(ctrPrices, ctrPrice)
				experiencedAds = append(experiencedAds, ad)
			}
		}
	}

	selectNewAd := RandomGenerator.Float64() < *NewAdSelectionProbability

	var finalAd *common.AdWithMetrics
	if selectNewAd && len(newAds) > 0 {
		selectedAd := newAds[RandomGenerator.Intn(len(newAds))]
		finalAd = &selectedAd
	} else if len(experiencedAds) > 0 {
		totalScore := 0.0
		for _, ctrPrice := range ctrPrices {
			totalScore += ctrPrice
		}

		randomPoint := RandomGenerator.Float64() * totalScore
		currentSum := 0.0
		for i, ad := range experiencedAds {
			currentSum += ctrPrices[i]
			if randomPoint <= currentSum {
				finalAd = &ad
				break
			}
		}
	} else if len(newAds) > 0 {
		// Handle fall back
		selectedAd := newAds[RandomGenerator.Intn(len(newAds))]
		finalAd = &selectedAd
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
		fmt.Printf("\nFailed to get image\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	publisherID, err := strconv.ParseUint(pubID, 10, 64)
	if err != nil {
		fmt.Printf("\nFailed to get publisherID\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publisher ID"})
		return
	}

	var impression = common.UrlImpressionParameters{
		ID:         uuid.New(),
		Pid:        int(publisherID),
		AdId:       int(ad.Id),
		IsClicked:  false,
		LoadAdTime: time.Now(),
	}

	var click = common.UrlClickParameters{
		ID:           uuid.New(),
		Pid:          int(publisherID),
		AdId:         int(ad.Id),
		ImpressionID: impression.ID,
		ExpTime:      time.Now().Add(clickExpirationTime),
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
