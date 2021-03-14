package propose

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	proposeMain "training/internal/propose"
)

func Create(main proposeMain.ProposeUsecases) func(w http.ResponseWriter, h *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Println("Parsing propose", r.Body)
		propose, err := main.ParsePropose(r.Body)
		logrus.Println("Propose parsed", propose)
		if err != nil {
			errMessage := err.Error()
			http.Error(w,  errMessage, http.StatusInternalServerError)
			return
		}
		logrus.Println("Validating propose", propose)
		err = main.Validate(propose)
		if err != nil {
			errMessage := err.Error()
			http.Error(w,  errMessage, http.StatusBadRequest)
			return
		}
		logrus.Println("Checking if user has propose", propose)
		err = main.CheckIfUserAlreadyHasPropose(propose.Cpf)
		if err != nil {
			errMessage := err.Error()
			http.Error(w,  errMessage, http.StatusUnprocessableEntity)
			return
		}
		clientResponse, err := main.CheckIfProposeIsAvailable(propose)
		if err != nil {
			errMessage := err.Error()
			http.Error(w,  errMessage, http.StatusInternalServerError)
			return
		}
		go main.GenerateCard(clientResponse)
		proposeDomain := clientResponse.ToPropose(propose)
		err = main.Create(proposeDomain)
		if err != nil {
			errMessage := err.Error()
			http.Error(w,  errMessage, http.StatusInternalServerError)
		}
		ticker := time.NewTicker(500 * time.Millisecond)
		done := make(chan bool)
		go func() {
			for {
				select{
					case <- done:
						return
					case <-ticker.C :
						err := main.CheckIfCardWasGenerated(proposeDomain)
						if err != nil {
							http.Error(w,  err.Error(), http.StatusInternalServerError)
							ticker.Stop()
						}
				}
			}
		}()
		time.Sleep(10 * time.Second)
		done <- true
	}
}
