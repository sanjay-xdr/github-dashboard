package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sanjay-xdr/github-dashboard/backend/internals/models"
)

func GetPullRequestData(w http.ResponseWriter, r *http.Request) {

	// var task models.Data;

	// task=&{
	// 	Name:"Sanjay"
	// 	Age:13
	// }
	// // err := json.NewDecoder(r.Body).Decode((&task))
	// // if err != nil {
	// // 	http.Error(w, err.Error(), http.StatusBadRequest)
	// // 	return
	// // }

	dummyData := &models.PRStatus{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dummyData)

}
