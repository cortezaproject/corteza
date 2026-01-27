package governor

func (g *governor) PauseAll() {
	if g.paused.CompareAndSwap(false, true) {
		g.mu.Lock()
		g.gates.globalPause = newGate()
		g.mu.Unlock()
	}
}

func (g *governor) ResumeAll() {
	if g.paused.CompareAndSwap(true, false) {
		g.mu.Lock()
		g.gates.globalPause.open()
		g.mu.Unlock()
	}
}
