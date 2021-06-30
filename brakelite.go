package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

const TimeScale = time.Second

var breakDurs = [...]int{5, 10, 15, 20, 25, 30}
var pauseDurs = [...]int{30, 60, 90, 120}
var breakStartDurIdx = 1
var breakMsgs = [...]string{
	"Take a quick break, sit up straight.",
	"Grab a sip of water, and unbend your spine.",
}

type DurOpt struct {
	menu *systray.MenuItem
	dur  int
}

var durOptAggCh chan *DurOpt
var pauseOptAggCh chan *DurOpt
var notifTimerStopCh chan struct{}

var durOpts []DurOpt
var pauseOpts []DurOpt
var menuStatus *systray.MenuItem
var timeUntilNextNotif int
var timeUntilUnpause int = 0

func main() {
	durOptAggCh = make(chan *DurOpt)
	pauseOptAggCh = make(chan *DurOpt)
	notifTimerStopCh = make(chan struct{})
	systray.Run(onReady, onExit)
}

func genNotifMsg() string {
	return breakMsgs[rand.Intn(len(breakMsgs))]
}

func showNotif() {
	if err := beeep.Notify("Brakelite!", genNotifMsg(), "icon.ico"); err != nil {
		panic(err)
	}
}

func resetNotifTimer(dopt *DurOpt, isFirst bool) {
	if !isFirst {
		notifTimerStopCh <- struct{}{}
	}

	for _, dopt2 := range durOpts {
		if *dopt == dopt2 {
			dopt2.menu.Check()
		} else {
			dopt2.menu.Uncheck()
		}
	}

	if timeUntilUnpause <= 0 {
		timeUntilNextNotif = dopt.dur
		menuStatus.SetTitle(fmt.Sprintf("Next: %d mins", timeUntilNextNotif))
	}

	go func() {
		for {
			select {
			case <-time.After(TimeScale):
				if timeUntilUnpause > 0 {
					timeUntilUnpause--
					menuStatus.SetTitle(fmt.Sprintf("Paused: %d mins", timeUntilUnpause))
				}
				if timeUntilUnpause <= 0 {
					timeUntilNextNotif--
					if timeUntilNextNotif <= 0 {
						showNotif()
						timeUntilNextNotif = dopt.dur
					}
					menuStatus.SetTitle(fmt.Sprintf("Next: %d mins", timeUntilNextNotif))
				}
			case <-notifTimerStopCh:
				return
			}
		}
	}()
}

func addDurOpt(root *systray.MenuItem, dur int) {
	menuText := fmt.Sprintf("%d Mins", dur)
	menu := root.AddSubMenuItemCheckbox(menuText, menuText, false)
	dopt := DurOpt{menu, dur}
	durOpts = append(durOpts, dopt)

	go func() {
		for {
			select {
			case <-menu.ClickedCh:
				durOptAggCh <- &dopt
			}
		}
	}()
}

func addPauseOpt(root *systray.MenuItem, dur int) {
	menuText := fmt.Sprintf("%d Mins", dur)
	menu := root.AddSubMenuItem(menuText, menuText)
	dopt := DurOpt{menu, dur}
	pauseOpts = append(pauseOpts, dopt)

	go func() {
		for {
			select {
			case <-menu.ClickedCh:
				pauseOptAggCh <- &dopt
			}
		}
	}()
}

func onReady() {
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetTitle("Brakelite")
	systray.SetTooltip("Lightweight. Minimal. Effective.")

	menuStatus = systray.AddMenuItem("N/A", "Status text")
	menuStatus.Disable()

	menuDurs := systray.AddMenuItem("Durations", "Notification duration settings")
	for _, dur := range breakDurs {
		addDurOpt(menuDurs, dur)
	}
	resetNotifTimer(&durOpts[breakStartDurIdx], true)

	menuPauseDurs := systray.AddMenuItem("Pause", "Pause brakelite for some time")
	menuCancelPause := menuPauseDurs.AddSubMenuItem("Cancel Pause", "Cancel any current pause")
	for _, dur := range pauseDurs {
		addPauseOpt(menuPauseDurs, dur)
	}

	menuQuit := systray.AddMenuItem("Quit", "Quit Brakelite")

	go func() {
		for {
			select {
			case dopt := <-durOptAggCh:
				resetNotifTimer(dopt, false)
			case dopt := <-pauseOptAggCh:
				timeUntilUnpause = dopt.dur
				timeUntilNextNotif = 0
				menuStatus.SetTitle(fmt.Sprintf("Paused: %d mins", timeUntilUnpause))
			case <-menuCancelPause.ClickedCh:
				timeUntilUnpause = 0
			case <-menuQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	notifTimerStopCh <- struct{}{}
	close(durOptAggCh)
	close(pauseOptAggCh)
	close(notifTimerStopCh)
}

// TODO: Load from icon.ico file, of course
var iconData = []byte{
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x18, 0x18, 0x02, 0x00, 0x01, 0x00,
	0x01, 0x00, 0xf0, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00,
	0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xff, 0x00, 0xfc, 0x00, 0x3f, 0x00, 0xf8, 0x00, 0x1f, 0x00, 0xf0, 0x7e,
	0x0f, 0x00, 0xe0, 0xff, 0x07, 0x00, 0xc1, 0xff, 0x83, 0x00, 0x83, 0xff,
	0xc1, 0x00, 0x87, 0xff, 0xe1, 0x00, 0x8f, 0xff, 0xf1, 0x00, 0x9f, 0xff,
	0xf9, 0x00, 0x9f, 0xff, 0xf9, 0x00, 0x9e, 0x00, 0x79, 0x00, 0x9e, 0x00,
	0x79, 0x00, 0x9f, 0xff, 0xf9, 0x00, 0x9f, 0xff, 0xf9, 0x00, 0x8f, 0xff,
	0xf1, 0x00, 0x87, 0xff, 0xe1, 0x00, 0x83, 0xff, 0xc1, 0x00, 0xc1, 0xff,
	0x83, 0x00, 0xe0, 0xff, 0x07, 0x00, 0xf0, 0x7e, 0x0f, 0x00, 0xf8, 0x00,
	0x1f, 0x00, 0xfc, 0x00, 0x3f, 0x00, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff,
	0xff, 0x00, 0xfc, 0x00, 0x3f, 0x00, 0xf8, 0x00, 0x1f, 0x00, 0xf0, 0x7e,
	0x0f, 0x00, 0xe0, 0xff, 0x07, 0x00, 0xc1, 0xff, 0x83, 0x00, 0x83, 0xff,
	0xc1, 0x00, 0x87, 0xff, 0xe1, 0x00, 0x8f, 0xff, 0xf1, 0x00, 0x9f, 0xff,
	0xf9, 0x00, 0x9f, 0xff, 0xf9, 0x00, 0x9e, 0x00, 0x79, 0x00, 0x9e, 0x00,
	0x79, 0x00, 0x9f, 0xff, 0xf9, 0x00, 0x9f, 0xff, 0xf9, 0x00, 0x8f, 0xff,
	0xf1, 0x00, 0x87, 0xff, 0xe1, 0x00, 0x83, 0xff, 0xc1, 0x00, 0xc1, 0xff,
	0x83, 0x00, 0xe0, 0xff, 0x07, 0x00, 0xf0, 0x7e, 0x0f, 0x00, 0xf8, 0x00,
	0x1f, 0x00, 0xfc, 0x00, 0x3f, 0x00, 0xff, 0xff, 0xff, 0x00,
}
