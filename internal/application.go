package internal

import (
	"currency_converter/currency_converter/utils"
	"os"
)

func CreateCurrFile() {
	file, err := os.OpenFile(utils.FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	utils.Logger.Info("file created")
	err = file.Close()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
}
