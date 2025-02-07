// Package plugin_gif
// @Description:
package plugin_gif

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"image"
	"strings"

	"github.com/guonaihong/gout"
	"github.com/huoxue1/gg"
	log "github.com/sirupsen/logrus"

	"github.com/huoxue1/leafBot"
	"github.com/huoxue1/leafBot/message"
)

func init() {
	LuXun()

}

// LuXun
/**
 * @Description:
 * example
 */
func LuXun() {
	plugin := leafBot.NewPlugin("鲁迅说")
	plugin.SetHelp(map[string]string{"发送鲁迅说即可获取结果": ""})
	plugin.OnStartWith("鲁迅说").SetPluginName("鲁迅说").SetWeight(10).AddHandle(func(event leafBot.Event, bot leafBot.Api, state *leafBot.State) {
		text := event.GetPlainText()
		data := strings.TrimLeft(text, "鲁迅说")
		if len(data) == 0 {
			event.Send(message.Text("你想让鲁迅说点什么呢？"))
			event1, err := event.GetOneEvent(func(event1 leafBot.Event, bot1 leafBot.Api, state1 *leafBot.State) bool {
				if event1.UserId == event.UserId && event1.GroupId == event.GroupId {
					return true
				}
				return false
			})
			if err != nil {
				return
			}
			data = event1.GetPlainText()
		}

		img, err := getImage(data)
		if err != nil {
			event.Send(message.Text("鲁迅说出错了" + err.Error()))
			return
		}
		event.Send(message.Image("base64://" + base64.StdEncoding.EncodeToString(img)))
	})
}

func getImage(text string) ([]byte, error) {
	var result []byte
	buffer := bytes.NewBuffer(result)
	var data []byte
	err := gout.GET("https://codechina.csdn.net/m15082717021/image/-/raw/main/202109090936718.jpeg").BindBody(&data).Do()
	if err != nil {
		return nil, err
	}
	decode, s, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	log.Infoln(s)
	context := gg.NewContextForImage(decode)
	err = context.LoadFontFromBytes(leafBot.GetFont(), 30)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	context.SetHexColor("FFFFFF")
	context.DrawString("——鲁迅", 320, 440)
	log.Infoln(len(text))
	context.DrawString(text, 240-float64(len(text)*5), 370)
	err = context.EncodePNG(buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), err
}
