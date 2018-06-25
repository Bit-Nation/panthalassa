package api

import (
	"time"

	pb "github.com/Bit-Nation/panthalassa/api/pb"
)

type DApp struct {
	api *API
}

// request to show a modal
func (a *DApp) ShowModal(title, layout string) error {

	// send request
	_, err := a.api.request(&pb.Request{
		ShowModal: &pb.Request_ShowModal{
			Title:  title,
			Layout: layout,
		},
	}, time.Second*20)

	return err

}
