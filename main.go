package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"os"
)

func main() {
	// 查找中文字体的存在
	fontPath, _ := findfont.Find("WenQuanDengKuanWeiMiHei.ttf")
	_ = os.Setenv("FYNE_FONT", fontPath) // 设置带中文的字体
	fmt.Println(os.Getenv("FYNE_FONT"))

	a := app.New()
	w := a.NewWindow("中文测试")
	w.Resize(fyne.NewSize(512, 512))
	w.SetContent(widget.NewLabel("中文测试一手"))

	w.Show()
	a.Run()
	_ = os.Unsetenv("FYNE_FONT")
}
