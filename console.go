package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func ConsoleWidget() []Widget {
	var searchDir *walk.LineEdit
	var pngCheck, jpegCheck, bmpCheck, heicCheck *walk.CheckBox

	return []Widget{
		Label{
			Text: "Search Directory: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				LineEdit{
					AssignTo: &searchDir,
					Text:     ConfigGet().SearchDir,
					OnEditingFinished: func() {
						dir := searchDir.Text()
						if dir != "" {
							stat, err := os.Stat(dir)
							if err != nil {
								ErrorBoxAction(mainWindow, "Search directory is not exist")
								searchDir.SetText("")
								SearchDirSave("")
								return
							}
							if !stat.IsDir() {
								ErrorBoxAction(mainWindow, "Search directory is not directory")
								searchDir.SetText("")
								SearchDirSave("")
								return
							}
						}
						SearchDirSave(dir)
					},
				},
				PushButton{
					MaxSize: Size{Width: 30},
					Text:    " ... ",
					OnClicked: func() {
						dlgDir := new(walk.FileDialog)
						dlgDir.FilePath = ConfigGet().SearchDir
						dlgDir.Flags = win.OFN_EXPLORER
						dlgDir.Title = "Please select a folder as search directory"

						exist, err := dlgDir.ShowBrowseFolder(mainWindow)
						if err != nil {
							logs.Error(err.Error())
							return
						}
						if exist {
							logs.Info("select %s as search directory", dlgDir.FilePath)
							searchDir.SetText(dlgDir.FilePath)
							SearchDirSave(dlgDir.FilePath)
						}
					},
				},
			},
		},
		Label{
			Text: "Search Options: ",
		},
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				CheckBox{
					AssignTo: &pngCheck,
					Text:     IMG_PNG,
					Checked:  SelectGet(IMG_PNG),
					OnCheckedChanged: func() {
						SelectCheck(IMG_PNG, pngCheck.Checked())
					},
				},
				CheckBox{
					AssignTo: &jpegCheck,
					Text:     IMG_JPEG,
					Checked:  SelectGet(IMG_JPEG),
					OnCheckedChanged: func() {
						SelectCheck(IMG_JPEG, jpegCheck.Checked())
					},
				},
				CheckBox{
					AssignTo: &bmpCheck,
					Text:     IMG_BMP,
					Checked:  SelectGet(IMG_BMP),
					OnCheckedChanged: func() {
						SelectCheck(IMG_BMP, bmpCheck.Checked())
					},
				},
				CheckBox{
					AssignTo: &heicCheck,
					Text:     IMG_HEIC,
					Checked:  SelectGet(IMG_HEIC),
					OnCheckedChanged: func() {
						SelectCheck(IMG_HEIC, heicCheck.Checked())
					},
				},
			},
		},
	}
}
