package service

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Структуры для работы с XML-ответом
type KeyRateEnvelope struct {
	XMLName xml.Name    `xml:"Envelope"`
	Body    KeyRateBody `xml:"Body"`
}

type KeyRateBody struct {
	Response KeyRateResponse `xml:"KeyRateXMLResponse"`
}

type KeyRateResponse struct {
	Result KeyRateResult `xml:"KeyRateXMLResult"`
}

type KeyRateResult struct {
	Rows []KeyRateRows `xml:"KeyRate"`
}

type KeyRateRows struct {
	KeyRates []KeyRates `xml:"KR"`
}

type KeyRates struct {
	Date string `xml:"DT" json:"date"`
	Rate string `xml:"Rate" json:"rate"`
}

type CbrService interface {
	GetLastKeyRate() (KeyRates, error)
}

type cbrService struct{}

// NewCbrService returns a new instance of CbrService
func NewCbrService() CbrService {
	return &cbrService{}
}

// Функция для получения ключевой ставки
func getKeyRate(fromDate, toDate string) ([]KeyRates, error) {
	request := GetKeyRateXMLRequest{
		Xmlns:    "http://web.cbr.ru/",
		FromDate: fromDate,
		ToDate:   toDate,
	}

	rawXmlData := SoapCall("https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx", request)

	var data KeyRateEnvelope
	if err := xml.Unmarshal([]byte(rawXmlData), &data); err != nil {
		return nil, err
	}

	return data.Body.Response.Result.Rows[0].KeyRates, nil
}

// Структура для SOAP-запроса
type GetKeyRateXMLRequest struct {
	XMLName  xml.Name `xml:"KeyRateXML"`
	Xmlns    string   `xml:"xmlns,attr"`
	FromDate string   `xml:"fromDate"`
	ToDate   string   `xml:"ToDate"`
}

// Функция для отправки SOAP-запроса
func SoapCall(service string, request GetKeyRateXMLRequest) string {
	var root = struct {
		XMLName xml.Name `xml:"soap12:Envelope"`
		Xsi     string   `xml:"xmlns:xsi,attr"`
		Xsd     string   `xml:"xmlns:xsd,attr"`
		Soap12  string   `xml:"xmlns:soap12,attr"`
		Body    struct {
			XMLName xml.Name             `xml:"soap12:Body"`
			Request GetKeyRateXMLRequest `xml:"KeyRateXML"`
		}
	}{
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:    "http://www.w3.org/2001/XMLSchema",
		Soap12: "http://www.w3.org/2003/05/soap-envelope",
	}
	root.Body.Request = request

	out, _ := xml.MarshalIndent(&root, " ", "  ")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	response, err := client.Post(service, "application/soap+xml", bytes.NewBufferString(string(out)))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	return string(content)
}

// Функция для получения последней доступной ключевой ставки
func (c *cbrService) GetLastKeyRate() (KeyRates, error) {
	today := time.Now()
	for daysBack := 0; daysBack < 30; daysBack++ { // Проверяем последние 30 дней
		dateToCheck := today.AddDate(0, 0, -daysBack)
		keyRates, err := getKeyRate(dateToCheck.Format("2006-01-02"), dateToCheck.Format("2006-01-02"))
		if err != nil {
			return KeyRates{}, err
		}
		if len(keyRates) > 0 {
			return keyRates[0], nil // Возвращаем последнюю найденную ключевую ставку
		}
	}
	return KeyRates{}, fmt.Errorf("нет доступных ключевых ставок за последние 30 дней")
}
