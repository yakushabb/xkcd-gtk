package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strconv"
)

// GotoDialog holds a gtk dialog that asks the user for a comic number
// to jump to.
type GotoDialog struct {
	parent *Window
	dialog *gtk.Dialog
	entry  *gtk.Entry
}

// NewGotoDialog creates and returns a GotoDialog with parent set as
// dialog's parent.
func NewGotoDialog(parent *Window) (*GotoDialog, error) {
	var err error
	gt := new(GotoDialog)
	gt.parent = parent
	gt.dialog, err = gtk.DialogNew()
	if err != nil {
		return nil, err
	}
	gt.dialog.SetTransientFor(parent.win)
	gt.dialog.SetTitle("Go to comic number...")
	gt.dialog.SetResizable(false)
	gt.dialog.SetDestroyWithParent(true)
	gt.dialog.Connect("delete-event", gt.Destroy)
	gt.dialog.Connect("response", gt.Response)

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 12)
	if err != nil {
		return nil, err
	}
	box.SetMarginStart(12)
	box.SetMarginEnd(12)

	icon, err := gtk.ImageNewFromIconName("dialog-question", gtk.ICON_SIZE_DIALOG)
	if err != nil {
		return nil, err
	}
	icon.SetMarginTop(12)
	icon.SetMarginBottom(12)
	icon.SetVAlign(gtk.ALIGN_CENTER)
	box.Add(icon)

	label, err := gtk.LabelNew("Go to")
	if err != nil {
		return nil, err
	}
	box.Add(label)
	gt.entry, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	gt.entry.SetActivatesDefault(true)
	gt.entry.SetPlaceholderText("Comic #")
	box.Add(gt.entry)
	_, err = gt.dialog.AddButton("Cancel", 0)
	if err != nil {
		return nil, err
	}
	submit, err := gt.dialog.AddButton("Go", 1)
	if err != nil {
		return nil, err
	}
	submitStyle, err := submit.GetStyleContext()
	if err != nil {
		return nil, err
	}
	submitStyle.AddClass("suggested-action")
	submit.SetCanDefault(true)
	submit.GrabDefault()
	box.ShowAll()

	contentArea, err := gt.dialog.GetContentArea()
	if err != nil {
		return nil, err
	}
	contentArea.Add(box)

	return gt, nil
}

// Present is a wrapper around gt.dialog.Present()
func (gt *GotoDialog) Present() {
	gt.dialog.Present()
}

// Destroy removes our references to the dialog so the garbage collector
// can take care of it.
func (gt *GotoDialog) Destroy() {
	gt.entry = nil
	gt.dialog = nil
	gt.parent.gotoDialog = nil
	gt.parent = nil
}

// Response is called when a `response` signal is received by the
// dialog.
func (gt *GotoDialog) Response(dialog *gtk.Dialog, responseID gtk.ResponseType) {
	defer dialog.Close()
	if responseID == 1 {
		input, err := gt.entry.GetText()
		if err != nil {
			log.Print(err)
			return
		}
		number, err := strconv.Atoi(input)
		if err != nil {
			log.Print(err)
			return
		}
		gt.parent.SetComic(number)
	}
}
