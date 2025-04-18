package internal

import (
	"encoding/json"
	"io"
	"net/http"
)

// Структура валюты с кодом и курсом
//type Valute struct {
//	Code string
//	Rate float64
//}

//type Currency interface {
//	GetCurrency()
//	GetRate()
//}

// Структура с курсами валют на момент даты запроса для json
type Rates struct {
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// Запрос к API для получения актуальных курсов валют
func RatesRequest() (Rates, error) {
	var result Rates
	resp, err := http.Get("https://www.cbr-xml-daily.ru/latest.js")

	if err != nil {
		return Rates{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Rates{}, err
	}

	err = json.Unmarshal(body, &result)

	if err != nil {
		return Rates{}, err
	}

	return result, nil
}

// Получение курса валюты по коду
func (r *Rates) GetRate(code string) float64 {
	for key, rate := range r.Rates {
		if key == code {
			return rate
		}
	}
	return 0
}

// Получение всех доступных валют
func (r *Rates) GetCurrency() map[string]float64 {
	return r.Rates
}

//func ParseValute(rates []string, currency Currency) []Valute {
//	var parsed []Valute
//	for i, rate := range rates {
//		parsed = append(parsed, Valute{rate, currency.Rates[rate]})
//	}
//}
