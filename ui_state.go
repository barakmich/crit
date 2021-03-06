package crit

type uiStateFunc func() error

type UIState struct {
	theme          *Theme
	review         *reviewState
	selectedCommit *reviewCommit
	listeners      []uiStateFunc
	app            *app
	fileSet        *fileSet
}

func (ui *UIState) registerChange(f uiStateFunc) {
	ui.listeners = append(ui.listeners, f)
}

func (ui *UIState) update(f func(ui *UIState) error) error {
	err := f(ui)
	for _, x := range ui.listeners {
		err = x()
		if err != nil {
			return err
		}
	}
	return nil
}
