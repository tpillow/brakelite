package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

var breakDurs = [...]int{5, 10, 15, 20}
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
var notifTimerStopCh chan struct{}

var durOpts []DurOpt

func main() {
	durOptAggCh = make(chan *DurOpt)
	notifTimerStopCh = make(chan struct{})
	systray.Run(onReady, onExit)
}

func genNotifMsg() string {
	return breakMsgs[rand.Intn(len(breakMsgs))]
}

func showNotif() {
	// TODO: Icon
	if err := beeep.Notify("Brakelite! Stop!", genNotifMsg(), ""); err != nil {
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

	go func() {
		for {
			select {
			case <-time.After(time.Duration(dopt.dur) * time.Second):
				showNotif()
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

func onReady() {
	// TODO: systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Brakelite")
	systray.SetTooltip("Lightweight. Minimal. Effective.")

	menuDurs := systray.AddMenuItem("Durations", "Notification duration settings")
	for _, dur := range breakDurs {
		addDurOpt(menuDurs, dur)
	}
	resetNotifTimer(&durOpts[breakStartDurIdx], true)

	menuQuit := systray.AddMenuItem("Quit", "Quit Brakelite")

	go func() {
		for {
			select {
			case dopt := <-durOptAggCh:
				resetNotifTimer(dopt, false)
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
	close(notifTimerStopCh)
}
