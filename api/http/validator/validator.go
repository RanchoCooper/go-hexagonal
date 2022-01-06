package validator

import (
    "strings"

    "github.com/gin-gonic/gin"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
)

/**
 * @author Rancho
 * @date 2022/1/6
 */

type ValidError struct {
    Key     string
    Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
    return v.Message
}

func (v *ValidErrors) Errors() []string {
    var errs []string
    for _, err := range *v {
        errs = append(errs, err.Error())
    }
    return errs
}

func (v *ValidErrors) Error() string {
    return strings.Join(v.Errors(), ",")
}

func BindAndValid(c *gin.Context, v interface{}, binder func(interface{}) error) (bool, ValidErrors) {
    var errs ValidErrors
    err := binder(v)
    if err != nil {
        v := c.Value("trans")
        trans, _ := v.(ut.Translator)
        verrs, ok := err.(validator.ValidationErrors)
        if !ok {
            return false, errs
        }

        for key, value := range verrs.Translate(trans) {
            errs = append(errs, &ValidError{
                Key:     key,
                Message: value,
            })
        }
        return false, errs
    }

    return true, nil
}
