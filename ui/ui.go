package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/matheusd/nonota"
	"github.com/rivo/tview"
)

const codeWidth = 45

type NonotaUI struct {
	board    *nonota.Board
	user     *nonota.User
	app      *tview.Application
	refTime  time.Time
	filename string

	tree        *tview.TreeView
	rootNode    *tview.TreeNode
	detailPages *tview.Pages
	root        *tview.Grid

	editor       *Editor
	timeForm     *tview.Form
	gridTaskForm *tview.Grid
	lastWork     *nonota.Work
	confirmWork  *nonota.Work
	statusBar    *tview.TextView
	treeNodes    map[interface{}]*tview.TreeNode
}

func New(board *nonota.Board, filename string, refTime time.Time) *NonotaUI {

	rootNode := tview.NewTreeNode("Board").SetSelectable(true).SetReference(board)
	tree := tview.NewTreeView().SetRoot(rootNode).SetCurrentNode(rootNode)
	tree.SetTitle("Lists").SetBorder(true)
	tree.SetVimBindingsEnabled(true)

	timeForm := tview.NewForm()
	timeForm.
		SetBorder(true).
		SetTitle("Confirm Time Input")

	editor := NewEditor()

	gridTaskForm := tview.NewGrid().
		SetRows(-5, -1).
		AddItem(editor.GetPrimitive(), 0, 0, 1, 1, 0, 0, false)

	detailPages := tview.NewPages().
		AddPage("timeConfirm", timeForm, true, true).
		AddPage("editor", gridTaskForm, true, true)

	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	root := tview.NewGrid().
		SetRows(-1, 1).
		SetColumns(-8, -4).
		AddItem(tree, 0, 0, 1, 1, 0, 0, true).
		AddItem(detailPages, 0, 1, 1, 1, 0, 0, false).
		AddItem(statusBar, 1, 0, 1, 2, 0, 0, false)

		//AddItem(detailPages, codeWidth, 1, false)

	app := tview.NewApplication().SetRoot(root, true)
	ui := &NonotaUI{
		filename:     filename,
		board:        board,
		refTime:      refTime,
		user:         &nonota.User{},
		app:          app,
		rootNode:     rootNode,
		tree:         tree,
		detailPages:  detailPages,
		root:         root,
		timeForm:     timeForm,
		editor:       editor,
		gridTaskForm: gridTaskForm,
		statusBar:    statusBar,
		treeNodes:    make(map[interface{}]*tview.TreeNode),
	}

	timeForm.
		AddInputField("Duration", "", 0, nil, nil).
		AddInputField("Note", "", 0, nil, nil)

	editor.CancelFunc = func() {
		ui.treeNodeSelected(tree.GetCurrentNode())
		ui.app.SetFocus(tree)
	}

	editor.AcceptFunc = func() {
		ui.updateCurrentNode()
		ui.app.SetFocus(tree)
	}

	ui.setInputCapture()
	ui.recreateLists()
	tree.SetChangedFunc(ui.treeNodeSelected)
	rootNode.ExpandAll()

	go func() {
		for {
			time.Sleep(time.Second)
			ui.app.QueueUpdateDraw(ui.perSecondUpdate)
		}
	}()

	return ui
}

func (ui *NonotaUI) save() {
	err := nonota.BoardToFile(ui.filename, ui.board)
	if err != nil {
		fmt.Println("EEEERRR", err)
	}
}

func (ui *NonotaUI) setInputCapture() {
	ui.tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currNode := ui.tree.GetCurrentNode()

		switch r := currNode.GetReference().(type) {
		case *nonota.List:
			switch {
			case event.Rune() == 'J':
				ui.board.MoveDown(r)
				ui.save()
			case event.Rune() == 'K':
				ui.board.MoveUp(r)
				ui.save()
			case event.Rune() == 'i':
				ui.app.SetFocus(ui.editor.GetPrimitive())
			case event.Rune() == 'a':
				ui.board.AppendNewTask(r)
			default:
				return event
			}
		case *nonota.Task:
			switch {
			case event.Rune() == 'J':
				ui.board.MoveTaskDown(r)
				ui.save()
			case event.Rune() == 'K':
				ui.board.MoveTaskUp(r)
				ui.save()
			case event.Rune() == 'i':
				ui.app.SetFocus(ui.editor.GetPrimitive())
			case event.Rune() == ' ':
				ui.lastWork = ui.user.ToggleWorkOnTask(r)
			case event.Key() == tcell.KeyEnter:
				ui.confirmToStopWork(r)
			default:
				return event
			}
		case *nonota.Board:
			switch {
			case event.Rune() == 'a':
				ui.board.AppendNewList()
			default:
				return event
			}
		default:
			return event
		}
		ui.app.QueueUpdateDraw(ui.recreateLists)

		return nil
	})
}

func (ui *NonotaUI) perSecondUpdate() {
	var txt string

	if ui.lastWork != nil {
		workTime := ui.lastWork.CurrentDuration().Round(time.Second)
		txt += " ⌚" + workTime.String() + " " + ui.lastWork.Task.Title + "\t"
	}

	now := ui.refTime
	dayTotal := ui.board.TotalTime(nonota.StartOfDay(now), nonota.EndOfDay(now))
	weekTotal := ui.board.TotalTime(nonota.StartOfWeek(now), nonota.EndOfWeek(now))
	billTotal := ui.board.TotalTime(nonota.StartOfBilling(now), nonota.EndOfBilling(now))

	txt += fmt.Sprintf("⌚ day %s week %s bill %s", dayTotal, weekTotal, billTotal)

	if ui.statusBar.GetText(false) != txt {
		ui.statusBar.Clear()
		ui.statusBar.SetText(txt)
	}
}

func (ui *NonotaUI) confirmToStopWork(task *nonota.Task) {
	work := ui.user.WorkForTask(task)
	if work == nil {
		work = nonota.NewWork(task)
	}

	ui.confirmWork = work
	workTime := work.CurrentDuration().Round(time.Minute) + time.Minute
	workTimeStr := workTime.String()
	workTimeStr = workTimeStr[:len(workTimeStr)-2]

	for ui.timeForm.GetButtonCount() > 0 {
		ui.timeForm.RemoveButton(ui.timeForm.GetButtonCount() - 1)
	}

	ui.detailPages.SwitchToPage("timeConfirm")
	ui.app.SetFocus(ui.timeForm)
	fldDuration := ui.timeForm.GetFormItem(0).(*tview.InputField)
	fldDuration.SetText(workTimeStr)

	fldNote := ui.timeForm.GetFormItem(1).(*tview.InputField)
	fldNote.SetText("")

	ui.timeForm.
		AddButton("Confirm", func() {
			fldDuration := ui.timeForm.GetFormItem(0).(*tview.InputField)
			fldNote := ui.timeForm.GetFormItem(1).(*tview.InputField)
			workTime, err := time.ParseDuration(fldDuration.GetText())
			if err != nil {
				ui.user.ExcludeWork(ui.confirmWork)
			} else {
				ui.confirmWork.SetNote(fldNote.GetText())
				ui.confirmWork.AdjustWorkDuration(workTime)
				ui.user.StopWork(ui.confirmWork)
			}
			if ui.lastWork == ui.confirmWork {
				ui.lastWork = nil
			}
			ui.detailPages.SwitchToPage("editor")
			ui.recreateLists()
			ui.confirmWork = nil
			ui.app.SetFocus(ui.tree)
			ui.save()
		}).
		AddButton("Cancel", func() {
			ui.detailPages.SwitchToPage("editor")
			ui.confirmWork = nil
			simulEvent := tcell.NewEventKey(tcell.KeyTAB, 0, tcell.ModNone)
			ui.timeForm.InputHandler()(simulEvent, func(tview.Primitive) {})
			ui.app.SetFocus(ui.tree)
		})
}

func (ui *NonotaUI) updateCurrentNode() {
	currNode := ui.tree.GetCurrentNode()

	text := ui.editor.GetText()
	firstLine := text
	descr := ""
	if idxln := strings.IndexRune(text, '\n'); idxln > 0 {
		firstLine = text[:idxln]
		descr = text[idxln+1:]
	}

	switch r := currNode.GetReference().(type) {
	case *nonota.List:
		r.Title = firstLine
	case *nonota.Task:
		r.Title = firstLine
		r.Description = descr
	default:
		return
	}

	ui.save()

	ui.recreateLists()
}

func (ui *NonotaUI) recreateLists() {
	startTime := nonota.StartOfBilling(ui.refTime)
	endTime := nonota.EndOfBilling(ui.refTime)

	children := make([]*tview.TreeNode, len(ui.board.Lists))
	var selNode *tview.TreeNode

	currNode := ui.tree.GetCurrentNode()
	selItem := currNode.GetReference()

	for i, l := range ui.board.Lists {
		totalTime := l.TotalTime(startTime, endTime).Round(time.Minute)
		text := l.Title
		if totalTime > 0 {
			totalTimeStr := totalTime.String()
			text += " ⌚" + totalTimeStr[:len(totalTimeStr)-2] // possible to do, since rounded.
		}

		n, has := ui.treeNodes[l]
		if !has {
			n = tview.NewTreeNode(text).SetSelectable(true).SetReference(l)
			n.SetColor(tcell.ColorGreen)
			ui.treeNodes[l] = n
		} else {
			n.SetText(text)
		}

		if l == selItem {
			selNode = n
		}

		children[i] = n

		listNodes := make([]*tview.TreeNode, len(l.Tasks))

		for j, t := range l.Tasks {
			text := t.Title
			totalTime := t.TotalTime(startTime, endTime)
			if totalTime > 0 {
				text += " ⌚" + totalTime.Round(time.Second).String()
			}

			tn, has := ui.treeNodes[t]
			if !has {
				tn = tview.NewTreeNode(text).SetSelectable(true).SetReference(t)
				ui.treeNodes[t] = tn
			} else {
				tn.SetText(text)
			}

			listNodes[j] = tn

			if ui.lastWork != nil && t == ui.lastWork.Task {
				tn.SetColor(tcell.ColorYellow)
			}

			if t == selItem {
				selNode = tn
			}
		}

		n.SetChildren(listNodes)
	}

	ui.rootNode.SetChildren(children)
	if selNode != nil {
		ui.tree.SetCurrentNode(selNode)
	}
}

func (ui *NonotaUI) treeNodeSelected(node *tview.TreeNode) {
	switch r := node.GetReference().(type) {
	case *nonota.Task:
		ui.detailPages.SwitchToPage("editor")
		fullText := r.Title
		if r.Description != "" {
			fullText += "\n\n" + strings.TrimLeft(r.Description, "\n")
		}
		ui.editor.SetText(fullText)
	case *nonota.List:
		_ = r
		ui.detailPages.SwitchToPage("editor")
		ui.editor.SetText(r.Title)
	}
}

func (ui *NonotaUI) Run() error {
	return ui.app.Run()
}
