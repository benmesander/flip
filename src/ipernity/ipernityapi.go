package ipernity

import (
	"encoding/json"
	"errors"
	"strconv"
)

func Call_user_get(userid string) (Userget, error) {
	var (
		parms pslice
	)

	parms = append(parms, Parameter{"api_key", apikey}, Parameter{"auth_token", token})
	if userid != "" {
		parms = append(parms, Parameter{"user_id", userid})
	}
	f, err := CallApiMethod(parms, "user.get")
	if err != nil {
		return Userget{}, err
	}
	jsonresult := &Userget{}
	json.Unmarshal(f, &jsonresult)

	if jsonresult.Api.Status != "ok" {
		return *jsonresult, errors.New("Error getting user data: " + jsonresult.Api.Code + " " + jsonresult.Api.Message)
	}

	return *jsonresult, nil
}

func Call_doc_getList(userid string, page int, extra string) (Docgetlist, error) {
	var (
		parms pslice
	)

	parms = append(parms, Parameter{"api_key", apikey}, Parameter{"auth_token", token}, Parameter{"per_page", "100"}, Parameter{"page", strconv.Itoa(page)})
	if userid != "" {
		parms = append(parms, Parameter{"user_id", userid})
	}
	if extra != "" {
		parms = append(parms, Parameter{"extra", extra})
	}
	f, err := CallApiMethod(parms, "doc.getList")
	if err != nil {
		return Docgetlist{}, err
	}
	jsonresult := &Docgetlist{}
	json.Unmarshal(f, &jsonresult)

	if jsonresult.Api.Status != "ok" {
		return *jsonresult, errors.New("Error getting document list: " + jsonresult.Api.Code + " " + jsonresult.Api.Message)
	}

	return *jsonresult, nil
}

func Call_doc_get(doc_id string, extra string) (Docget, error) {
	var (
		parms pslice
	)

	parms = append(parms, Parameter{"api_key", apikey}, Parameter{"auth_token", token}, Parameter{"doc_id", doc_id})
	if extra != "" {
		parms = append(parms, Parameter{"extra", extra})
	}
	f, err := CallApiMethod(parms, "doc.get")
	if err != nil {
		return Docget{}, err
	}
	jsonresult := &Docget{}
	json.Unmarshal(f, &jsonresult)

	if jsonresult.Api.Status != "ok" {
		return *jsonresult, errors.New("Error getting document: " + jsonresult.Api.Code + " " + jsonresult.Api.Message)
	}

	return *jsonresult, nil
}

func Call_doc_getContainers(doc_id string) (Docgetcontainers, error) {
	var (
		parms pslice
	)

	parms = append(parms, Parameter{"api_key", apikey}, Parameter{"auth_token", token}, Parameter{"doc_id", doc_id})
	f, err := CallApiMethod(parms, "doc.getContainers")
	if err != nil {
		return Docgetcontainers{}, err
	}
	jsonresult := &Docgetcontainers{}
	json.Unmarshal(f, &jsonresult)

	if jsonresult.Api.Status != "ok" {
		return *jsonresult, errors.New("Error getting document containers: " + jsonresult.Api.Code + " " + jsonresult.Api.Message)
	}

	return *jsonresult, nil
}

