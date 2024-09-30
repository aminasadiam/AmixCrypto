package cmd

import (
	"AmixCrypto/internal"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

var (
	AuthUsernames []string
	Intervals     []string
	startMessage  = &telebot.Message{Text: "<b>Let's get started!</b>🪙\nHello! please send me your token name for prediction"}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatalln("Couldn't Load .env File")
	}
	AuthUsernames = strings.Split(os.Getenv("AUTH_USERNAMES"), ",")
	Intervals = strings.Split(os.Getenv("INTERVALS"), ",")
}

func Execute(token string) error {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		return err
	}

	bot.Handle("/start", func(c telebot.Context) error {
		for _, username := range AuthUsernames {
			if username == c.Message().Sender.Username {
				return c.Send(startMessage.Text, &telebot.SendOptions{
					ParseMode: "HTML",
				})
			}
		}
		return c.Send(fmt.Sprintf("Hello, %s!, you dont have access for work with this bot. sorry!", c.Sender().Username))
	})

	bot.Handle("/crypto", func(c telebot.Context) error {
		argv := c.Args()
		if len(argv) > 2 {
			return c.Reply("لطفا فقط اسم یک ارز را بدهید")
		} else if len(argv) == 0 {
			return c.Reply("لطفا نام ارز و تایم فریم را نیز ارسال کنید")
		} else if len(argv) == 1 {
			symbol := strings.ToUpper(argv[0]) + "USDT"
			exist, err := internal.CheckSymbolExist(symbol)
			if err != nil {
				logrus.Error(err)
				return err
			}
			if exist {
				return c.Reply(fmt.Sprintf("%v Price is : %v", symbol, internal.GetNowPrice(symbol)))
			} else {
				return c.Send(fmt.Sprintf("%s! Not Found", symbol))
			}
		} else {
			for _, username := range AuthUsernames {
				if username == c.Message().Sender.Username {
					symbol := strings.ToUpper(argv[0]) + "USDT"
					interval := strings.ToLower(argv[1])
					exist, err := internal.CheckSymbolExist(symbol)
					if err != nil {
						logrus.Error(err)
						return err
					}
					if exist {
						c.Send(fmt.Sprintf("در حال دریافت اطلاعات جهت تحلیل و پیشبینی ارز %v...", symbol))
						go func(symbol string) {
							err = internal.GetSymbolHistory(symbol, interval)
							if err != nil {
								logrus.Error(err)
							}
						}(symbol)
						time.Sleep(5 * time.Second)
						c.Send(fmt.Sprintf("در حال تحلیل ارز %v...", symbol))
						result, err := exec.Command("python", "TrainModel/train.py", symbol, interval).Output()
						if err != nil {
							logrus.Fatalln(err)
						}
						time.Sleep(2 * time.Second)
						return c.Send(fmt.Sprintf("نتیجه پیشبینی %v -> %v:\n\n\n%v", symbol, interval, string(result)))
					} else {
						return c.Send(fmt.Sprintf("%s! Not Found", symbol))
					}
				}
			}
			return c.Send(fmt.Sprintf("Hello, %s!, you dont have access for work with this bot. sorry!", c.Sender().Username))
		}
	})

	// bot.Handle("/forex", func(c telebot.Context) error {
	// 	argv := c.Args()
	// 	if len(argv) > 1 {
	// 		return c.Reply("لطفا فقط اسم یک ارز را بدهید")
	// 	} else if len(argv) == 0 {
	// 		return c.Reply("لطفا نام ارز را نیز ارسال کنید")
	// 	} else {
	// 		for _, username := range AuthUsernames {
	// 			if username == c.Message().Sender.Username {
	// 				symbol := strings.ToUpper(argv[0])
	// 				exist, err := internal.CheckForexExist(symbol)
	// 				if err != nil {
	// 					logrus.Error(err)
	// 					return err
	// 				}
	// 				if exist {
	// 					c.Send(fmt.Sprintf("در حال دریافت اطلاعات جهت تحلیل و پیشبینی ارز %v...", symbol))
	// 					go func(symbol string) {
	// 						err = internal.ForexHistoricalData(symbol)
	// 						if err != nil {
	// 							logrus.Error(err)
	// 						}
	// 					}(symbol)
	// 					time.Sleep(5 * time.Second)
	// 					c.Send(fmt.Sprintf("در حال تحلیل ارز %v...", symbol))
	// 					result, err := exec.Command("python", "TrainModel/forex.py", symbol).Output()
	// 					if err != nil {
	// 						logrus.Fatalln(err)
	// 					}
	// 					time.Sleep(2 * time.Second)
	// 					return c.Send(fmt.Sprintf("نتیجه پیشبینی فارکس ارز %v:\n\n\n%v", symbol, string(result)))
	// 				} else {
	// 					return c.Send(fmt.Sprintf("%s! Not Found", symbol))
	// 				}
	// 			}
	// 		}
	// 		return c.Send(fmt.Sprintf("Hello, %s!, you dont have access for work with this bot. sorry!", c.Sender().Username))
	// 	}
	// })

	bot.Start()
	return nil
}
