package propose

import (
	"net/http"
	proposeMain "training/internal/propose"
)

func Create(main proposeMain.ProposeUsecases) func(w http.ResponseWriter, h *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		propose, err := main.ParsePropose(r.Body)
		if err != nil {
			errMessage := err.Error()
			http.Error(w,  errMessage, http.StatusInternalServerError)
		}
		err = main.Validate(propose)
		if err != nil {
			errMessage := err.Error()
			http.Error(w,  errMessage, http.StatusBadRequest)
		}
		main.Create(propose)
	}
}
