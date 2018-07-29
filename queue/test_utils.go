package queue

type testProcessor struct {
	processorType string
	validJob      func(j Job) error
	process       func(j Job) error
}

func (p *testProcessor) Type() string {
	return p.processorType
}

func (p *testProcessor) ValidJob(j Job) error {
	return p.validJob(j)
}

func (p *testProcessor) Process(j Job) error {
	return p.process(j)
}

type testStorage struct {
	persistJob func(j Job) error
	deleteJob  func(j string) error
}

func (s *testStorage) PersistJob(j Job) error {
	return s.persistJob(j)
}

func (s *testStorage) DeleteJob(id string) error {
	return s.deleteJob(id)
}
