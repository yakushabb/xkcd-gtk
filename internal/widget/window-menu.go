package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type WindowMenu struct {
	*gtk.MenuButton

	popover    *gtk.Popover
	popoverBox *gtk.Box

	zoomBox        *ZoomBox
	darkModeSwitch *DarkModeSwitch // may be nil
}

var _ Widget = &WindowMenu{}

func NewWindowMenu(accels *gtk.AccelGroup, prefersAppMenu bool, darkModeGetter func() bool, darkModeSetter func(bool)) (*WindowMenu, error) {
	super, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	wm := &WindowMenu{
		MenuButton: super,
	}

	wm.SetTooltipText(l("Window menu"))
	wm.AddAccelerator("activate", accels, gdk.KEY_F10, 0, gtk.ACCEL_VISIBLE)

	wm.popover, err = gtk.PopoverNew(wm)
	if err != nil {
		return nil, err
	}
	wm.SetPopover(wm.popover)
	wm.SetUsePopover(true)

	wm.popoverBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	wm.popover.Add(wm.popoverBox)
	defer wm.popoverBox.ShowAll()

	addMenuSeparator := func() error {
		sep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		if err != nil {
			return err
		}
		wm.popoverBox.PackStart(sep, false, true, style.PaddingPopoverCompact/2)
		return nil
	}

	addMenuEntry := func(label, action string, external bool) error {
		mb, err := gtk.ModelButtonNew()
		if err != nil {
			return err
		}
		mb.SetActionName(action)
		if external {
			label = urlLabel(label)
		}
		mb.SetLabel(label)
		mbl, err := mb.GetChild()
		if err != nil {
			return err
		}
		mbl.ToWidget().SetHAlign(gtk.ALIGN_START)
		wm.popoverBox.PackStart(mb, false, true, 0)
		return nil
	}

	// Zoom section.
	wm.zoomBox, err = NewZoomBox()
	if err != nil {
		return nil, err
	}
	wm.zoomBox.SetMarginBottom(style.PaddingPopoverCompact / 2)
	wm.popoverBox.Add(wm.zoomBox)

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	// Comic properties section.
	err = addMenuEntry(l("Open link"), "win.open-link", true)
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("Explain"), "win.explain", true)
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("Properties"), "win.show-properties", false)
	if err != nil {
		return nil, err
	}

	// If the desktop environment will show an app menu, then we do not need
	// to add the app menu contents to the window menu.
	if prefersAppMenu {
		return wm, nil
	}

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	err = addMenuEntry(l("New window"), "app.new-window", false)
	if err != nil {
		return nil, err
	}

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	wm.darkModeSwitch, err = NewDarkModeSwitch(darkModeGetter, darkModeSetter)
	if err != nil {
		return nil, err
	}
	wm.popoverBox.PackStart(wm.darkModeSwitch, false, true, 0)

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	err = addMenuEntry(l("What If?"), "app.open-what-if", true)
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("xkcd blog"), "app.open-blog", true)
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("xkcd store"), "app.open-store", true)
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("About xkcd"), "app.open-about-xkcd", true)
	if err != nil {
		return nil, err
	}

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	err = addMenuEntry(l("Keyboard shortcuts"), "app.show-shortcuts", false)
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("About"), "app.show-about", false)
	if err != nil {
		return nil, err
	}

	return wm, nil
}

func (wm *WindowMenu) Dispose() {
	if wm == nil {
		return
	}

	wm.MenuButton = nil

	wm.popover = nil
	wm.popoverBox = nil
	wm.zoomBox.Dispose()
	wm.zoomBox = nil
	wm.darkModeSwitch.Dispose()
	wm.darkModeSwitch = nil
}

func (wm *WindowMenu) SetCompact(compact bool) {
	if compact {
		wm.popoverBox.SetMarginTop(style.PaddingPopoverCompact)
		wm.popoverBox.SetMarginBottom(style.PaddingPopoverCompact)
		wm.popoverBox.SetMarginStart(0)
		wm.popoverBox.SetMarginEnd(0)
	} else {
		wm.popoverBox.SetMarginTop(style.PaddingPopover)
		wm.popoverBox.SetMarginBottom(style.PaddingPopover)
		wm.popoverBox.SetMarginStart(style.PaddingPopover)
		wm.popoverBox.SetMarginEnd(style.PaddingPopover)
	}
	wm.zoomBox.SetCompact(compact)
	wm.darkModeSwitch.SetCompact(compact)
}
