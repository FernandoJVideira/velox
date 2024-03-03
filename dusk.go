package velox

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"time"
)

func (v *Velox) TakeScreenShot(pageURl, testName string, w, h float64) {
	page := v.FetchPage(pageURl)

	img, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  w,
			Height: h,
			Scale:  1,
		},
		FromSurface: true,
	})
	if err != nil {
		v.ErrorLog.Println(err)
	}
	fileName := time.Now().Format("02-01-2006-15-04-05.000000")
	_ = utils.OutputFile(fmt.Sprintf("%s/screenshots/%s-%s.png", v.RootPath, testName, fileName), img)
}

func (v *Velox) FetchPage(url string) *rod.Page {
	page := rod.New().MustConnect().MustIgnoreCertErrors(true).MustPage(url).MustWaitLoad()
	return page
}

func (v *Velox) SelectElementById(page *rod.Page, id string) *rod.Element {
	return page.MustElementByJS(fmt.Sprintf("document.getElementById('%s')", id))
}
