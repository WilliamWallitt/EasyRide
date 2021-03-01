package Error_Management

import (
	"github.com/go-playground/locales/eu"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"net/http"
)

type Error struct {
	ResponseCode int
	Errors []string
}

type Auth struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type Driver struct {
	Rate int `validate:"required"`
}

type Trip struct {
	Origin string `validate:"required"`
	Destination string `validate:"required"`

}


func FormValidationHandler(model interface{}) (*interface{}, Error) {

	// create new validator instance
	v := validator.New()
	// create new Error struct instance
	response := Error{}

	if err := v.Struct(model); err != nil {
		translator := eu.New()
		uni := ut.New(translator, translator)
		trans, _ := uni.GetTranslator("en")

		if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
			response.ResponseCode = http.StatusInternalServerError
			return nil, response
		}

		for _, err := range err.(validator.ValidationErrors) {
			response.Errors = append(response.Errors, err.Translate(trans))
		}
		response.ResponseCode = http.StatusBadRequest
		return nil, response
	}

	response.ResponseCode = http.StatusOK
	return &model, response

}


