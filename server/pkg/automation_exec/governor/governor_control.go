package governor

func (g *governor) PauseAll() {
	if g.paused.CompareAndSwap(false, true) {
		g.mux.Lock()
		g.gates.globalPause = newGate()
		g.mux.Unlock()
	}
}

func (g *governor) ResumeAll() {
	if g.paused.CompareAndSwap(true, false) {
		g.mux.Lock()
		g.gates.globalPause.open()
		g.mux.Unlock()
	}
}
