package widget

import (
	"github.com/gotk3/gotk3/gtk"
)

type BookmarkCheckButton struct {
	*gtk.CheckButton

	bookmarked    func() bool
	setBookmarked func(bool)
}

var _ Widget = &BookmarkCheckButton{}

func NewBookmarkCheckButton(bookmarkedGetter func() bool, bookmarkedSetter func(bool)) (*BookmarkCheckButton, error) {
	super, err := gtk.CheckButtonNew()
	if err != nil {
		return nil, err
	}
	bcb := &BookmarkCheckButton{
		CheckButton: super,

		bookmarked:    bookmarkedGetter,
		setBookmarked: bookmarkedSetter,
	}

	bcb.SetLabel(l("Bookmark this comic"))
	bcb.Connect("toggled", bcb.CheckStateChanged)
	bcb.Update()

	return bcb, nil
}

func (bcb *BookmarkCheckButton) Dispose() {
	bcb.CheckButton = nil
	bcb.bookmarked = nil
	bcb.setBookmarked = nil
}

func (bcb *BookmarkCheckButton) CheckStateChanged() {
	active := bcb.GetActive()
	// Avoid calling bcb.setBookmarked when this signal might have been emitted by
	// bcb.Update.
	if active == bcb.bookmarked() {
		return
	}
	bcb.setBookmarked(active)
}

func (bcb *BookmarkCheckButton) Update() {
	isBookmarked := bcb.bookmarked()
	// Avoid calling bcb.SetActive when this signal might have been emitted by
	// bcb.CheckStateChanged.
	if isBookmarked == bcb.GetActive() {
		return
	}
	bcb.SetActive(isBookmarked)
}
