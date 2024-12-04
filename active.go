package main

import (
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var search, cancel *walk.PushButton

func ButtonEnable() {
	ProcessUpdate(1000)
	time.Sleep(time.Millisecond * 500)
	search.SetEnabled(true)
	cancel.SetEnabled(true)
	ProcessUpdate(0)
}

func ButtonDisable() {
	search.SetEnabled(false)
	cancel.SetEnabled(false)
}

func ActiveWidget() []Widget {
	return []Widget{
		PushButton{
			AssignTo: &search,
			Text:     "Search",
			OnClicked: func() {
				ButtonDisable()
				go func() {
					SearchFileActive()
					ButtonEnable()
				}()
			},
		},
		HSpacer{},
		PushButton{
			AssignTo: &cancel,
			Text:     "Cancel",
			OnClicked: func() {
				CloseWindows()
			},
		},
	}
}
