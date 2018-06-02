package main

import (
  "fmt"
  "strings"
  "time"

  "github.com/jroimartin/gocui"
)

func returnFocus(g *gocui.Gui, v *gocui.View) error {
  previousView, err := g.View(state.PreviousView)
  if err != nil {
    panic(err)
  }
  return switchFocus(g, v, previousView)
}

func switchFocus(g *gocui.Gui, oldView, newView *gocui.View) error {
  if oldView != nil {
    oldView.Highlight = false
    devLog("setting previous view to:", oldView.Name())
    state.PreviousView = oldView.Name()
  }
  newView.Highlight = true
  devLog(newView.Name())
  if _, err := g.SetCurrentView(newView.Name()); err != nil {
    return err
  }
  g.Cursor = newView.Editable
  return newLineFocused(g, newView)
}

func getItemPosition(v *gocui.View) int {
  _, cy := v.Cursor()
  _, oy := v.Origin()
  return oy + cy
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
  if v == nil {
    return nil
  }

  ox, oy := v.Origin()
  cx, cy := v.Cursor()
  if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
    if err := v.SetOrigin(ox, oy-1); err != nil {
      return err
    }
  }

  newLineFocused(g, v)
  return nil
}

func resetOrigin(v *gocui.View) error {
  if err := v.SetCursor(0, 0); err != nil {
    return err
  }
  return v.SetOrigin(0, 0)
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
  if v != nil {
    cx, cy := v.Cursor()
    ox, oy := v.Origin()
    if cy+oy >= len(v.BufferLines())-2 {
      return nil
    }
    if err := v.SetCursor(cx, cy+1); err != nil {
      if err := v.SetOrigin(ox, oy+1); err != nil {
        return err
      }
    }
  }

  newLineFocused(g, v)
  return nil
}

// if the cursor down past the last item, move it up one
func correctCursor(v *gocui.View) error {
  cx, cy := v.Cursor()
  _, oy := v.Origin()
  lineCount := len(v.BufferLines()) - 2
  if cy >= lineCount-oy {
    return v.SetCursor(cx, lineCount-oy)
  }
  return nil
}

func renderString(g *gocui.Gui, viewName, s string) error {
  g.Update(func(*gocui.Gui) error {
    timeStart := time.Now()
    v, err := g.View(viewName)
    if err != nil {
      panic(err)
    }
    v.Clear()
    fmt.Fprint(v, s)
    v.Wrap = true
    devLog("render time: ", time.Now().Sub(timeStart))
    return nil
  })
  return nil
}

func splitLines(multilineString string) []string {
  if multilineString == "" || multilineString == "\n" {
    return make([]string, 0)
  }
  lines := strings.Split(multilineString, "\n")
  if lines[len(lines)-1] == "" {
    return lines[:len(lines)-1]
  }
  return lines
}