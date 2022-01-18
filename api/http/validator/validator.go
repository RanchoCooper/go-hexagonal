package validator

import (
    "reflect"
    "strings"

    "github.com/fatih/structtag"
    "github.com/gin-gonic/gin"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
)

/**
 * @author Rancho
 * @date 2022/1/6
 */

const MessageTagKey = "message"

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

func BindAndValid(c *gin.Context, obj interface{}, binder func(interface{}) error) (bool, ValidErrors) {
    var errs ValidErrors
    err := binder(obj)
    if err != nil {
        v := c.Value("trans")
        trans, _ := v.(ut.Translator)
        verrs, ok := err.(validator.ValidationErrors)
        if !ok {
            return false, errs
        }

        for key, value := range verrs.Translate(trans) {
            validError := &ValidError{
                Key:     key,
                Message: value,
            }

            // get message tag, and replace valid Error.Message with message from tag
            tmpKey := strings.Split(key, ".")
            fieldName := tmpKey[len(tmpKey)-1]
            t := reflect.TypeOf(obj)
            k := t.Kind()
            if k == reflect.Ptr {
                t = t.Elem()
            }
            field, exists := t.FieldByName(fieldName)
            var tag reflect.StructTag
            if exists {
                tag = field.Tag
            }
            if tag != "" {
                tags, _ := structtag.Parse(string(tag))
                messageTag, _ := tags.Get(MessageTagKey)
                if messageTag != nil && messageTag.Name != "" {
                    validError.Message = messageTag.Name
                }
            }

            errs = append(errs, validError)
        }
        return false, errs
    }

    return true, nil
}
