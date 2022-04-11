package helper

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// BindAndValid struct bind and valid control
func BindAndValid(g *gin.Context, req interface{}) (int, interface{}) {

	if err := g.ShouldBind(&req); err != nil {
		return http.StatusBadRequest, err.Error()
	}

	valid := validation.Validation{}
	check, err := valid.Valid(req)
	if err != nil {
		return http.StatusInternalServerError, "Server error." + err.Error()
	}

	if !check {
		return http.StatusBadRequest, valid.Errors
	}
	return http.StatusOK, nil
}

// Valid only valid control
func Valid(g *gin.Context, req interface{}) (int, interface{}) {

	valid := validation.Validation{}
	check, err := valid.Valid(req)
	if err != nil {
		return http.StatusInternalServerError, "Server error." + err.Error()
	}

	if !check {
		return http.StatusBadRequest, valid.Errors
	}
	return http.StatusOK, nil
}
