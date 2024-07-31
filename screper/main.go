package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/liuzl/tokenizer"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

type PublisherToken struct {
	ID     uint     `gorm:"primaryKey;autoIncrement"`
	Name   string   `gorm:"column:name;unique;not null"`
	Url    string   `gorm:"column:url;unique;not null"`
	Tokens []*Token `gorm:"many2many:publisher_tokens;"`
}

type Publisher struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"column:name;unique;not null"`
	Url         string `gorm:"column:url;unique;not null"`
	Credit      int    `gorm:"column:credit"`
	Clicks      int    `gorm:"column:clicks"`
	Impressions int    `gorm:"column:impressions"`
}

type Token struct {
	ID         uint              `gorm:"primaryKey;autoIncrement"`
	Value      string            `gorm:"column:value;unique;not null"`
	Publishers []*PublisherToken `gorm:"many2many:publisher_tokens;"`
}

func htmlTokenizer(htmlText string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		log.Printf("failed to parse HTML: %v", err)
		return nil
	}
	text := doc.Text()
	tokens := tokenizer.Tokenize(text)
	log.Printf("tokenized %d tokens from text", len(tokens))
	return tokens
}

func fetchPageText(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code error for URL %s: %d %s", url, resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML for URL %s: %v", url, err)
	}
	text := doc.Text()
	return text, nil
}

func ScrapData() {
	var publishers []Publisher
	if err := DB.Find(&publishers).Error; err != nil {
		log.Fatalf("failed to fetch publishers: %v", err)
	}

	for _, publisher := range publishers {
		var publisherToken PublisherToken
		if err := DB.Where("id = ?", publisher.ID).First(&publisherToken).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				publisherToken = PublisherToken{
					ID:   publisher.ID,
					Name: publisher.Name,
					Url:  publisher.Url,
				}
				if err := DB.Create(&publisherToken).Error; err != nil {
					log.Printf("failed to create PublisherToken for publisher %d: %v", publisher.ID, err)
					continue
				}
				log.Printf("created PublisherToken for publisher %d", publisher.ID)
			} else {
				log.Printf("failed to fetch PublisherToken for publisher %d: %v", publisher.ID, err)
				continue
			}
		}

		text, err := fetchPageText(publisher.Url)
		if err != nil {
			log.Printf("failed to fetch page text for URL %s: %v", publisher.Url, err)
			continue
		}

		tokenizedText := htmlTokenizer(text)

		for _, tokenValue := range tokenizedText {
			var token Token
			if err := DB.Where("value = ?", tokenValue).First(&token).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					token = Token{
						Value: tokenValue,
					}
					if err := DB.Create(&token).Error; err != nil {
						log.Printf("failed to create Token with value %s: %v", tokenValue, err)
						continue
					}
					log.Printf("created Token with value %s", tokenValue)
				} else {
					log.Printf("failed to fetch Token with value %s: %v", tokenValue, err)
					continue
				}
			}

			if err := DB.Model(&token).Association("Publishers").Append(&publisherToken); err != nil {
				log.Printf("failed to associate Token with PublisherToken %d: %v", publisherToken.ID, err)
				continue
			}
			log.Printf("associated Token with value %s to PublisherToken %d", tokenValue, publisherToken.ID)
		}
	}
}

func ScrapRepetead() {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Starting new scraping cycle...")
				ScrapData()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {
	err := NewDatabase()
	if err != nil {
		log.Fatalf("cannot create database: %v", err)
	}

	err = DB.AutoMigrate(&Publisher{}, &Token{}, &PublisherToken{})
	if err != nil {
		log.Fatalf("cannot migrate schema: %v", err)
	}

	ScrapRepetead()
}
