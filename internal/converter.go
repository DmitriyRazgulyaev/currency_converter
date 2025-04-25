package internal

import (
	"currency_converter/currency_converter/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Структура с курсами валют на момент даты запроса для json
type Currency struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// Запрос к API для получения актуальных курсов валют
func RatesRequest() error {
	utils.Logger.Info("http request created")
	resp, err := http.Get("https://www.cbr-xml-daily.ru/latest.js")
	if err != nil {
		return err
	}

	if resp.Body == nil {
		panic("empty body response")
	}
	utils.Logger.Info("http request successfully earned")

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	utils.Logger.Info("response body read")

	var result Currency
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	utils.Logger.Info("json http request unmarshalled")

	err = WriteToFile([]Currency{result})
	if err != nil {
		return err
	}
	utils.Logger.Info("file written")

	return nil
}

// Получение курса валюты по коду на сегодняшнюю дату,
// если нет кода по такой дате возвращает 0, если нет такой даты - 0 и ошибку отсутствия даты
func GetTodayRates(currencies []Currency) (map[string]float64, error) {
	todayTime := time.Now()
	date := strings.Split(todayTime.String(), " ")[0]

	for _, curr := range currencies {
		if curr.Date == date {
			return curr.Rates, nil
		}
	}
	return nil, errors.New("date not found")
}

// Получение всех доступных валют
//func (r *Rates) GetCurrency() map[string]float64 {
//	return r.Rates
//}
//
//// Получение собственного курса выбранных валют
//func (r *Rates) GetUniqueCurrency(firstCode, secondCode string) float64 {
//	curr := r.Rates[secondCode] / r.Rates[firstCode]
//	return curr
//}

// Записывает данные в файл с курсами в виде время: {курсы}, возвращает количество записанных байт
func WriteToFile(data []Currency) error {

	var rates []Currency
	err := GetJson(&rates)
	if err != nil && err.Error() == "empty json" {
		MarshData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = os.WriteFile(utils.FILE, MarshData, 0666)
		if err != nil {
			return err
		}
		return nil
	}
	rates = append(rates, data...)
	updatedRates, err := json.Marshal(rates)
	if err != nil {
		return err
	}

	err = os.WriteFile(utils.FILE, updatedRates, 0666)
	if err != nil {
		return err
	}

	return nil
}

func GetJson(data *[]Currency) error {
	file, err := os.ReadFile(utils.FILE)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	if len(file) == 0 {
		return errors.New("empty json")
	}
	if !json.Valid(file) {
		utils.Logger.Fatal("invalid JSON format")

	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	utils.Logger.Info("getJson done successfully")
	return nil
}
