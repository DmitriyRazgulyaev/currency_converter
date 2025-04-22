package internal

import (
	"currency_converter/currency_converter/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

var logger = utils.Logger

// Структура с курсами валют на момент даты запроса для json
type Rates struct {
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// Запрос к API для получения актуальных курсов валют
func RatesRequest() (Rates, error) {
	utils.Logger.Info("http request created")
	resp, err := http.Get("https://www.cbr-xml-daily.ru/latest.js")
	if err != nil {
		return Rates{}, err
	}

	if resp.Body == nil {
		panic("empty body response")
	}
	utils.Logger.Info("http request successfully earned")

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Rates{}, err
	}
	utils.Logger.Info("response body read")

	var result Rates
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Rates{}, err
	}
	utils.Logger.Info("json unmarshalled")

	jsonRates, err := json.Marshal(map[string]map[string]float64{result.Date: result.Rates})
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	utils.Logger.Info("json marshalled")

	_, err = WriteToFile(jsonRates)
	if err != nil {
		return Rates{}, err
	}
	utils.Logger.Info("file written")

	return result, nil
}

// Получение курса валюты по коду
func (r *Rates) GetRate(code string) (float64, error) {
	for key, rate := range r.Rates {
		if key == code {
			return rate, nil
		}
	}
	return 0, errors.New("code absent")
}

// Получение всех доступных валют
func (r *Rates) GetCurrency() map[string]float64 {
	return r.Rates
}

// Получение собственного курса выбранных валют
func (r *Rates) GetUniqueCurrency(firstCode, secondCode string) float64 {
	curr := r.Rates[secondCode] / r.Rates[firstCode]
	return curr
}

// Записывает данные в файл с курсами в виде время: {курсы}, возвращает количество записанных байт
func WriteToFile(data []byte) (int, error) {
	file, err := os.OpenFile(utils.FILE, os.O_WRONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var n int
	n, err = file.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}
