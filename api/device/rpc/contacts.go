package rpc

type ContactListRequest struct{}

func (f *ContactListRequest) Type() string {
	return "CONTACT:LIST"
}

func (f *ContactListRequest) Data() (string, error) {
	return "", nil
}

func (f *ContactListRequest) Valid() error {
	return nil
}
