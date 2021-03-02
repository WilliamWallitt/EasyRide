package Error_Management

import (
	"github.com/go-playground/locales/eu"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"net/http"
)

// custom error struct
type Error struct {
	ResponseCode int
	Errors []string
}

// auth json request struct
type Auth struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

// driver json request struct
type Driver struct {
	Rate int `validate:"required"`
}

// trip json request struct
type Trip struct {
	Origin string `validate:"required"`
	Destination string `validate:"required"`

}

// generic form handler, that takes in a struct and makes sure that
// the json request is of the same format as the struct
func FormValidationHandler(model interface{}) (*interface{}, Error) {

	// create new validator instance
	v := validator.New()
	// create new Error struct instance
	response := Error{}
	// check that the current interface is of the correct format
	if err := v.Struct(model); err != nil {
		// custom error handling so the user can see what the problem is with their request
		translator := eu.New()
		uni := ut.New(translator, translator)
		trans, _ := uni.GetTranslator("en")
		if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
			response.ResponseCode = http.StatusInternalServerError
			return nil, response
		}
		// iterate over the request errors
		for _, err := range err.(validator.ValidationErrors) {
			// add them to the Error.Errors list
			response.Errors = append(response.Errors, err.Translate(trans))
		}
		// add the correct response code to the Error struct
		response.ResponseCode = http.StatusBadRequest
		// return the Error struct
		return nil, response
	}
	// if the request is good, return the interface
	// and make the response code StatusOK
	response.ResponseCode = http.StatusOK
	return &model, response

}


