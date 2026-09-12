package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"english-course-api/utils"

	"github.com/gin-gonic/gin"
)

func TestI18n_Translate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	utils.InitI18n()

	t.Run("Default English Translation", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Set("lang", "en")

		msg := utils.Translate(c, "student.created", nil)
		expected := "Student created successfully"
		if msg != expected {
			t.Errorf("expected '%s', got '%s'", expected, msg)
		}
	})

	t.Run("Indonesian Translation", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Set("lang", "id")

		msg := utils.Translate(c, "student.created", nil)
		expected := "Student berhasil dibuat"
		if msg != expected {
			t.Errorf("expected '%s', got '%s'", expected, msg)
		}
	})

	t.Run("Fallback for Non-Existent Key", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		c.Set("lang", "en")

		key := "non_existent.key"
		msg := utils.Translate(c, key, nil)
		if msg != key {
			t.Errorf("expected fallback '%s', got '%s'", key, msg)
		}
	})
}
