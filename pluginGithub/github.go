package pluginGithub

import (
	"fmt" //nolint:gci
	"regexp"
	"strings"

	"github.com/google/go-github/v35/github"
	"github.com/huoxue1/leafBot"
	"github.com/huoxue1/leafBot/message"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	client *github.Client
	ctx    context.Context
)

func init() {
	go PluginInit()
}

func githubInit() {
	if leafBot.DefaultConfig.Plugins.GithubToken == "" {
		log.Errorln("github_token not found")
	}

	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: leafBot.DefaultConfig.Plugins.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

func PluginInit() {
	githubInit()
	log.Infoln("已启用github插件，请确保你配置了github token")

	plugin := leafBot.NewPlugin("github解析")
	plugin.SetHelp(map[string]string{
		">github": "解析github仓库，例如：>github huoxue1/leafBot",
	})
	plugin.OnMessage("").
		SetPluginName("github url解析").
		SetWeight(10).
		AddRule(
			func(event leafBot.Event, bot leafBot.Api, state *leafBot.State) bool {
				compile := regexp.MustCompile(`https://github.com/([^\s]*)/([^\s]*)`)
				return compile.MatchString(event.Message.CQString())
			}).
		AddHandle(func(event leafBot.Event, bot leafBot.Api, state *leafBot.State) {
			log.Infoln("成功匹配")
			compile := regexp.MustCompile(`https://github.com/([^\s]*)/([^\s]*)`)
			datas := compile.FindStringSubmatch(event.Message.ExtractPlainText())
			if len(datas) != 3 {
				log.Errorln("正则匹配出现错误")
				return
			}
			msg, err := getResponseMsg(datas[1], datas[2])
			if err != nil {
				log.Errorln(err)
			}
			event.Send(msg)
		})

	plugin.OnCommand(">github").
		SetPluginName("github解析").
		SetWeight(10).
		SetBlock(false).
		SetCD("default", 0).
		AddHandle(func(event leafBot.Event, bot leafBot.Api, state *leafBot.State) {
			if len(state.Args) < 1 {
				event.Send("请输入你需要解析的仓库，例如\n>github huoxue1/leafBot")
				return
			}
			msg, err := getResponseMsg(strings.Split(state.Args[0], "/")[0], strings.Split(state.Args[0], "/")[1])
			if err != nil {
				event.Send("仓库获取失败" + err.Error())
			}
			event.Send([]message.MessageSegment{
				message.Text(msg),
				message.Image(
					fmt.Sprintf("https://opengraph.githubassets.com/0/%v/%v",
						strings.Split(state.Args[0], "/")[0],
						strings.Split(state.Args[0], "/")[1])).Add("cache", "0")})
		})
}

func getResponseMsg(owner, resp string) (string, error) {
	repository, _, err := client.Repositories.Get(ctx, owner, resp)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("%v\nDescription: %v\nStar/Fork/Issue: %d / %d / %d\nLanguage: %v\nLicense: %v\nLastPushed: %v\nJump:%v",
		repository.GetName(),
		repository.GetDescription(),
		repository.GetStargazersCount(),
		repository.GetForksCount(),
		repository.GetOpenIssuesCount(),
		repository.GetLanguage(),
		repository.GetLicense().GetName(),
		repository.GetPushedAt().Format("2006-01-02 15:04:05"),
		"https://github.com/"+owner+"/"+resp)
	return msg, err
}

func SearchResponse(ctx context.Context) ([]string, error) {
	var msgs []string

	repositories, _, err := client.Search.Repositories(ctx, "LeafBot", &github.SearchOptions{})
	if repositories == nil {
		return nil, err
	}
	for _, repository := range repositories.Repositories {
		msg := fmt.Sprintf("%v\nDescription: %v\nStar/Fork/Issue: %d / %d / %d\nLanguage: %v\nLicense: \nLastPushed: %v\nJump:%v",
			repository.GetName(),
			repository.GetDescription(),
			repository.GetStargazersCount(),
			repository.GetForksCount(),
			repository.GetOpenIssuesCount(),
			repository.GetLanguage(),
			//repository.GetLicense().String(),
			repository.GetPushedAt().Format("2006-01-02 15:04:05"),
			repository.GetURL(),
		)
		msgs = append(msgs, msg)
	}

	return msgs, err
}
