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
	return []Widget{
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text: "Search Directory: ",
				},
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
					MaxSize: Size{Width: 20},
					Text:    "...",
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
		Composite{
			Layout: HBox{MarginsZero: true},
			Children: []Widget{
				Label{
					Text: "Search Options: ",
				},
				CheckBox{
					Text:    "PNG",
					Checked: true,
				},
				CheckBox{
					Text:    "JPEG",
					Checked: true,
				},
				CheckBox{
					Text:    "BMP",
					Checked: true,
				},
				CheckBox{
					Text:    "HEIC",
					Checked: true,
				},
			},
		},
	}
}
