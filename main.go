package main

import (
	"encoding/json"
	"flag"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Anim struct {
	Type      string  `json:"__type__"`
	Name      string  `json:"_name"`
	ObjFlags  int     `json:"_objFlags"`
	Native    string  `json:"_native"`
	Duration  float64 `json:"_duration"`
	Sample    int     `json:"sample"`
	Speed     int     `json:"speed"`
	WrapMode  int     `json:"wrapMode"`
	CurveData struct {
		Comps struct {
			CcSprite struct {
				SpriteFrame []struct {
					Frame float64 `json:"frame"`
					Value struct {
						UUID string `json:"__uuid__"`
					} `json:"value"`
				} `json:"spriteFrame"`
			} `json:"cc.Sprite"`
		} `json:"comps"`
		Paths interface {
		} `json:"paths"`
	} `json:"curveData"`
	Events []interface{} `json:"events"`
}

var (
	frame = 1
	input = "./"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})

	startTime := time.Now().UnixNano()

	flag.IntVar(&frame, "frame", 1, "动画间隔帧")
	flag.StringVar(&input, "input", "./", "输入文件路径")
	flag.Parse()

	filepath.Walk(input, walkFunc)

	endTime := time.Now().UnixNano()
	log.Infof("总耗时:%v毫秒\n", (endTime-startTime)/1000000)
	time.Sleep(time.Millisecond * 1000)
}

func walkFunc(files string, info os.FileInfo, err error) error {
	if err != nil {
		log.Error(err)
		return err
	}
	_, fileName := filepath.Split(files)
	if path.Ext(files) == ".anim" {
		b, err := ioutil.ReadFile(filepath.FromSlash(input + "/" + fileName))
		if err != nil {
			log.Errorln(err.Error())
			return err
		}

		a := &Anim{}
		err = json.Unmarshal(b, a)
		if err != nil {
			log.Errorln(err.Error())
			return err
		}

		c := len(a.CurveData.Comps.CcSprite.SpriteFrame)
		var t = 1.0 / 60.0 * float64(frame)

		for i, _ := range a.CurveData.Comps.CcSprite.SpriteFrame {
			a.CurveData.Comps.CcSprite.SpriteFrame[i].Frame = t * float64(i)
		}

		for i, v := range a.CurveData.Comps.CcSprite.SpriteFrame {
			log.Infoln(i, v.Frame)
		}
		a.Duration = float64(c) * t

		if a.CurveData.Paths == nil {
			a.CurveData.Paths = make(map[string]interface{})
		}
		b, err = json.MarshalIndent(a, "", "  ")
		if err != nil {
			log.Errorln(err)
			return err
		}

		ioutil.WriteFile(fileName, b, 0777)
	}
	return nil
}
