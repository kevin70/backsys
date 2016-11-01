package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hoisie/redis"
	"github.com/urfave/cli"
)

type WxAccessToken struct {
	Errcode     int32  `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

var weixinCmd = &cli.Command{
	Name:    "weixin",
	Aliases: []string{"wx"},
	Usage:   "微信工具包",
	Subcommands: []*cli.Command{
		{
			Name:      "at",
			Usage:     "获取微信 access_token 并存放置指定的 redis 服务中，存储的 key 值等于 weixin.access_token.{{appid}}",
			UsageText: "backsys wx at --app-id [appid] --app-secret [appsecret]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "app-id",
					Usage: "REQUIRED 微信appid",
				},
				&cli.StringFlag{
					Name:  "app-secret",
					Usage: "REQUIRED 微信appsecret",
				},
				&cli.StringFlag{
					Name:  "redis-addr",
					Usage: "REQUIRED redis服务器地址，用于存储access_token",
				},
				&cli.StringFlag{
					Name:  "redis-pwd",
					Usage: "OPTIONAL redis服务器密码",
				},
				&cli.IntFlag{
					Name:  "redis-db-idx",
					Value: 0,
					Usage: "OPTIONAL redis数据库索引",
				},
			},
			Action: func(c *cli.Context) error {
				appId := c.String("app-id")
				if appId == "" {
					fmt.Println("缺少微信 [--app-id] 参数")
					return nil
				}

				appSecret := c.String("app-secret")
				if appSecret == "" {
					fmt.Println("缺少微信 [--app-secret] 参数")
					return nil
				}

				redisAddr := c.String("redis-addr")
				if redisAddr == "" {
					fmt.Println("缺少 redis 服务器地址")
					return nil
				}

				url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appId, appSecret)
				resp, err := http.Get(url)
				if err != nil {
					return cli.Exit(err.Error(), 2)
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return cli.Exit(err.Error(), 2)
				}

				var wxAccToken WxAccessToken
				err = json.Unmarshal(body, &wxAccToken)
				if err != nil {
					return cli.Exit(err.Error(), 2)
				}

				if wxAccToken.Errcode != 0 {
					return cli.Exit(fmt.Sprintf("获取微信[access_token]失败：%s", body), 2)
				}

				client := &redis.Client{
					Addr:        redisAddr,
					Db:          c.Int("redis-db-idx"),
					Password:    c.String("redis-pwd"),
					MaxPoolSize: 1,
				}

				key := fmt.Sprintf("weixin.access_token.%s", appId)
				err = client.Set(key, []byte(wxAccToken.AccessToken))
				if err != nil {
					return err
				}

				fmt.Println(string(body))
				return nil
			},
		},
	},
}

func AddWeixinCmd(app *cli.App) {
	app.Commands = append(app.Commands, weixinCmd)
}
