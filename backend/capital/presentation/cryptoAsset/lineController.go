package cryptoasset

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	linebot "github.com/line/line-bot-sdk-go/v7/linebot"
	share "github.com/ozaki-physics/raison-me/capital/infrastructure/share"
	usecase "github.com/ozaki-physics/raison-me/capital/usecase/cryptoAsset"
)

type LineController interface {
	// Get(w http.ResponseWriter, r *http.Request) error
	// Post(w http.ResponseWriter, r *http.Request) error
	// Put(w http.ResponseWriter, r *http.Request) error
	// Delete(w http.ResponseWriter, r *http.Request) error
	SoundReflection(w http.ResponseWriter, r *http.Request)
}
type lineController struct {
	credential share.CredentialLine
	caUsecase  usecase.CryptoAssetsUsecase
}

func CreateLineController(cr share.CredentialLine, u usecase.CryptoAssetsUsecase) LineController {
	return &lineController{cr, u}
}

func (linec *lineController) SoundReflection(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(linec.credential.Secret(), linec.credential.Token())
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(r)
	log.Printf("line request: %v\n", events)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Printf("line request message: %s\n", message.Text)
				msg := linec.mekeMessage(message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ozaki-physics:\n"+msg)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (linec *lineController) mekeMessage(symbol string) string {
	symbolUp := strings.ToUpper(symbol)
	gainPrice, err := linec.caUsecase.CoinGainPrice(symbolUp)
	if err != nil {
		log.Println(err)
	}

	var msg string
	msg += "symbol: " + symbol + "\n"
	msg += strconv.FormatFloat(gainPrice, 'f', -1, 64)
	return msg
}
