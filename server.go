package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Define the Go struct representing the XML data
type Alert struct {
	XMLName    xml.Name `xml:"urn:oasis:names:tc:emergency:cap:1.1 alert"`
	Identifier string   `xml:"identifier"`
	Sender     string   `xml:"sender"`
	Sent       string   `xml:"sent"`
	Status     string   `xml:"status"`
	MsgType    string   `xml:"msgType"`
	Scope      string   `xml:"scope"`
	Source     string   `xml:"source"`
	Info       Info     `xml:"info"`
}

type Info struct {
	Language     string `xml:"language"`
	Category     string `xml:"category"`
	Event        string `xml:"event"`
	ResponseType string `xml:"responseType"`
	Urgency      string `xml:"urgency"`
	Severity     string `xml:"severity"`
	Certainty    string `xml:"certainty"`
	Expires      string `xml:"expires"`
	SenderName   string `xml:"senderName"`
	Headline     string `xml:"headline"`
	Description  string `xml:"description"`
	Contact      string `xml:"contact"`
	Area         Area   `xml:"area"`
}

type Area struct {
	AreaDesc string `xml:"areaDesc"`
	Polygon  string `xml:"polygon"`
	GeoCode  string `xml:"geocode"`
}

func main() {
	http.HandleFunc("/EmergencyBroadcastRequest", handleRequest)
	http.HandleFunc("/notify", handleNotify)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	xmlData, err := ioutil.ReadAll(r.Body)
	taiwanTimezone, err := time.LoadLocation("Asia/Taipei")
	currentTime := time.Now().In(taiwanTimezone)
	formattedTime := currentTime.Format("2006-01-02 15:04:05.000000 MST")
	fmt.Printf("Received Emergency Broadcast Request from CBE at %s\n", formattedTime)
	if err != nil {
		http.Error(w, "Error parsing body data", http.StatusBadRequest)
		return
	}
	


	fmt.Printf("Received XML data\n")
	/*fmt.Printf("Received XML data: %s\n", string(xmlData))*/
	var alertData Alert
	if err := xml.Unmarshal(xmlData, &alertData); err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing XML data", http.StatusBadRequest)
		return
	}
	
	/*taiwanTimezone, err := time.LoadLocation("Asia/Taipei")
	currentTime := time.Now().In(taiwanTimezone)*/
	
	data := make(map[string]string)
	data["serialNumber"] = alertData.Identifier[len(alertData.Identifier)-3:]
	data["messageType"] = alertData.MsgType
	if alertData.Info.Language == "en-US" {
		data["dataCodingScheme"] = "01"
	}
	if alertData.Info.Language == "zh-TW" {
		data["dataCodingScheme"] = "48"
	}
	switch {
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1113"
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "1114"
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1115"
	case alertData.Info.Severity == "Extreme" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "1116"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1117"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Immediate" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "1118"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Observed":
		data["messageIdentifier"] = "1119"
	case alertData.Info.Severity == "Severe" && alertData.Info.Urgency == "Expected" && alertData.Info.Certainty == "Likely":
		data["messageIdentifier"] = "111A"
	default:
		data["messageIdentifier"] = "1112"
	}
	timeFormat := "2006-01-02 15:04:05.000000 MST"
	timeSent, err := time.Parse(timeFormat, alertData.Sent)
	data["warningMessageContents"] = timeSent.Format(timeFormat) + alertData.Info.Headline + "\n" + alertData.Info.Description + "\n" + alertData.Info.Area.AreaDesc
	data["timeSentFromCBE"] = alertData.Sent
	data["tac"] = alertData.Info.Area.GeoCode
	data["timeReceived"] = currentTime.Format(timeFormat)
	subscribe()
	transfer(data)

	fmt.Printf("Started distributing the warning message...\n")
	sendEmergencyBroadcastResponse(w)	
}

func sendEmergencyBroadcastResponse(w http.ResponseWriter) {
	responseData := `<response><status>Success</status><description>Emergency Broadcast Received</description></response>`
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseData))

	taiwanTimezone, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}
	currentTime := time.Now().In(taiwanTimezone)
	formattedTime := currentTime.Format("2006-01-02 15:04:05.000000 MST")

	fmt.Printf("Sent Emergency Broadcast Response to CBE at %s\n", formattedTime)

}


func handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received the request body"))
}