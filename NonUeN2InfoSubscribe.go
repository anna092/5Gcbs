package main

import (
	"context"
	"encoding/json"

	"github.com/free5gc/openapi/Namf_Communication"
	"github.com/free5gc/openapi/models"
)

func subscribe() {
	subscribe := models.NonUeN2InfoSubscriptionCreateData{}
	// Specify the URL you want to send the request to
	// Create the request body
	jsonString := []byte(`{
		"globalRanNodeList": [
		  {
			"gNbId": {
			  "bitLength": 24,
			  "gNBValue": ""
			},
			"plmnId": {
				"mnc": "",
				"mcc": ""
			},
			"n3IwfId": "",
			"ngeNbId": ""
		  }
		],
		"anTypeList":[

		],
		"n2InformationClass": "PWS",
<<<<<<< HEAD
		"n2NotifyCallbackUri": "192.168.56.102:8000/notify",
=======
		"n2NotifyCallbackUri": "127.0.0.1:8080/notify",
>>>>>>> b762fddca16447050b3f338d5e6b9cd6d78d7bd3
		"nfId": "",
		"supportedFeatures": ""
	  }`)
	json.Unmarshal(jsonString, &subscribe)
	namfConfiguration := Namf_Communication.NewConfiguration()
<<<<<<< HEAD
	//namfConfiguration.SetBasePath("http://127.0.0.18:8000")
	namfConfiguration.SetBasePath("http://192.168.56.102:8000")
=======
	namfConfiguration.SetBasePath("http://127.0.0.18:8000")
>>>>>>> b762fddca16447050b3f338d5e6b9cd6d78d7bd3
	apiClient := Namf_Communication.NewAPIClient(namfConfiguration)
	_, _, _ = apiClient.NonUEN2MessagesSubscriptionsCollectionDocumentApi.NonUeN2InfoSubscribe(context.TODO(), subscribe)
}
