package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func httpValidationErrorResponse(c *gin.Context, validationErrors map[string][]string) {
	c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": http.StatusBadRequest,
		"error": map[string]interface{}{
			"message": "Validation error",
			"errors":  validationErrors,
		},
	})

}

func httpBindingErrorResponse(c *gin.Context, err error) {
	var errMsg map[string][]string
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		errMsg = map[string][]string{
			ute.Field: {fmt.Sprintf("should be a %s", ute.Type)},
		}
	} else {
		errMsg = map[string][]string{
			"error": {err.Error()},
		}
	}

	var detailedMessages []string
	for field, errors := range errMsg {
		for _, e := range errors {
			message := fmt.Sprintf("%s %s", field, e)
			detailedMessages = append(detailedMessages, message)
		}
	}
	sort.Strings(detailedMessages)

	c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status": http.StatusBadRequest,
		"error": map[string]interface{}{
			"message":           "Validation error",
			"errors":            errMsg,
			"detailed_messages": detailedMessages,
		},
	})

}

func httpInternalServerErrorResponse(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"status": http.StatusInternalServerError,
		"error": map[string]interface{}{
			"message": errorMessage,
		},
	})
}

func httpUnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"status": http.StatusUnauthorized,
		"error": map[string]string{
			"message": "Unauthorized",
		},
	})
}

func httpOkResponse(c *gin.Context, payload map[string]interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"results": payload,
	})
}

func httpAuthenticationFailed(c *gin.Context, message string) {
	c.JSON(401, map[string]interface{}{
		"status": 401,
		"error": map[string]interface{}{
			"message": message,
		},
	})
}

func httpNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, map[string]interface{}{
		"status": http.StatusNotFound,
		"error": map[string]interface{}{
			"message": message,
		},
	})
}
