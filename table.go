package main

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type ImageSize struct {
	Width, Height int
}

func (s *ImageSize) Size() int {
	return s.Width * s.Height
}

type FileItem struct {
	Index       int
	File        string
	Size        ImageSize
	SimilarFile string
	SimilarSize ImageSize
	Similarity  float64
	Status      string

	checked bool
}

type FileModel struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items []*FileItem
}

func (n *FileModel) RowCount() int {
	return len(n.items)
}

func (n *FileModel) Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.File
	case 2:
		return fmt.Sprintf("%d*%d", item.Size.Width, item.Size.Height)
	case 3:
		return item.SimilarFile
	case 4:
		return fmt.Sprintf("%d*%d", item.SimilarSize.Width, item.SimilarSize.Height)
	case 5:
		return fmt.Sprintf("%.2f%%", item.Similarity)
	case 6:
		return item.Status
	}
	panic("unexpected col")
}

func (n *FileModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *FileModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *FileModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)
		case 1:
			return c(a.File < b.File)
		case 2:
			return c(a.Size.Size() < b.Size.Size())
		case 3:
			return c(a.SimilarFile < b.SimilarFile)
		case 4:
			return c(a.SimilarSize.Size() < b.SimilarSize.Size())
		case 5:
			return c(a.Similarity < b.Similarity)
		case 6:
			return c(a.Status < b.Status)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

const (
	STATUS_DONE      = "found"
	STATUS_READ_FAIL = "read fail"
	STATUS_HASH_FAIL = "hash fail"
)

var tableView *walk.TableView
var fileSimilarTable *FileModel

func init() {
	fileSimilarTable = new(FileModel)
	fileSimilarTable.items = make([]*FileItem, 0)
}

func SearchFileActive() {
	lt := fileSimilarTable

	lt.Lock()
	defer lt.Unlock()

	tableView.SetEnabled(false)
	defer tableView.SetEnabled(true)

	if ConfigGet().SearchDir == "" {
		ErrorBoxAction(mainWindow, "Please set the correct search directory!")
		return
	}

	stat, err := os.Stat(ConfigGet().SearchDir)
	if err != nil {
		ErrorBoxAction(mainWindow, "The search directory not exist!")
		return
	}

	if !stat.IsDir() {
		ErrorBoxAction(mainWindow, "The search directory is not directory!")
		return
	}

	lt.items = make([]*FileItem, 0)
	lt.PublishRowsReset()
	lt.Sort(lt.sortColumn, lt.sortOrder)

	fileList, err := ReadFileList(ConfigGet().SearchDir)
	if err != nil {
		ErrorBoxAction(mainWindow, fmt.Sprintf("Read the %s search directory fail, %s", ConfigGet().SearchDir, err.Error()))
		return
	}

	fileHashList := make([]*ImageInfo, 0)

	i := 0
	for index, file := range fileList {
		hashNew, err := ImageOpen(file, ConfigGet().SelectList)
		if err != nil {
			logs.Info(err.Error())
			continue
		}

		for _, hashOld := range fileHashList {
			similarity, err := ImageSimilarity(hashOld.hash, hashNew.hash)
			if err != nil {
				logs.Info(err.Error())
				continue
			}

			if similarity >= ConfigGet().Similarity {
				logs.Info("find the similarity %.2f%% image file, %s <-> %s", similarity, hashOld.file, hashNew.file)
				lt.items = append(lt.items, &FileItem{
					Index:       i,
					File:        hashOld.file,
					Size:        hashOld.size,
					SimilarFile: hashNew.file,
					SimilarSize: hashNew.size,
					Similarity:  similarity,
					Status:      STATUS_DONE})
				i++

				break
			}
		}

		fileHashList = append(fileHashList, hashNew)

		lt.PublishRowsReset()
		lt.Sort(lt.sortColumn, lt.sortOrder)

		ProcessUpdate(index * 1000 / len(fileList))
	}
}

func TableItemShow() {
	lt := fileSimilarTable

	lt.Lock()
	defer lt.Unlock()

	index := tableView.CurrentIndex()
	if index < len(lt.items) {
		item := lt.items[index]

		go OpenBrowserWeb(item.File)
		go OpenBrowserWeb(item.SimilarFile)
	}
}

func TableWidget() []Widget {
	return []Widget{
		Label{
			Text: "Image Similarity List:",
		},
		TableView{
			AssignTo:         &tableView,
			AlternatingRowBG: true,
			ColumnsOrderable: true,
			CheckBoxes:       false,
			OnItemActivated: func() {
				TableItemShow()
			},
			Columns: []TableViewColumn{
				{Title: "No", Width: 30},
				{Title: "File", Width: 250},
				{Title: "Size", Width: 120},
				{Title: "Similar File", Width: 250},
				{Title: "Size", Width: 120},
				{Title: "Similarity", Width: 80},
				{Title: "Status", Width: 60},
			},
			StyleCell: func(style *walk.CellStyle) {
				if style.Row()%2 == 0 {
					style.BackgroundColor = walk.RGB(248, 248, 255)
				} else {
					style.BackgroundColor = walk.RGB(220, 220, 220)
				}
			},
			Model: fileSimilarTable,
		},
	}
}
