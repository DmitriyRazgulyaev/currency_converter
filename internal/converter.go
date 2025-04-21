package internal

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

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

	if resp.Body == nil {
		panic("empty body response")
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
