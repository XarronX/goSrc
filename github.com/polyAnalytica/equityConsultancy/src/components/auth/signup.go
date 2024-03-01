package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/polyAnalytica/equityConsultancy/src/types"
)

func Signup(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		issue := types.NewIssue(err.Error())
		json.NewEncoder(writer).Encode(issue)
		return
	}

	_, err = types.NewClient(data)
	if err != nil {
		issue := types.NewIssue(err.Error())
		json.NewEncoder(writer).Encode(issue)
		return
	}

	response := types.NewResponse("success")
	json.NewEncoder(writer).Encode(response)
}

func Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		issue := types.NewIssue(err.Error())
		json.NewEncoder(writer).Encode(issue)
		return
	}

	info := map[string]string{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		issue := types.NewIssue(err.Error())
		json.NewEncoder(writer).Encode(issue)
		return
	}

	_, err = types.LoadClient(info["email"], info["passHash"])
	if err != nil {
		issue := types.NewIssue(err.Error())
		json.NewEncoder(writer).Encode(issue)
		return
	}

	response := types.NewResponse("success")
	json.NewEncoder(writer).Encode(response)
}

func checkIfUserExists(productId string) bool {
	// var product types.User
	// db.Instance.First(&product, productId)
	// if product.ID == 0 {
	// 	return false
	// }
	return true
}
