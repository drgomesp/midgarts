package opengl

type State struct {
	program     *Program
	bufferCount int
}

func (s *State) Program() *Program {
	return s.program
}
