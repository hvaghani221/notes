package e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gavv/httpexpect/v2"

	"notes/internal/model"
)

func TestPost(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  url,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	if err := setup(e); err != nil {
		t.Fatal(err)
	}

	e.POST("/api/auth/login").WithJSON(model.LogInDTO{
		Email:    users[0].Email,
		Password: "bad password",
	}).Expect().Status(http.StatusUnauthorized)

	r := e.POST("/api/auth/login").WithJSON(model.LogInDTO{
		Email:    users[0].Email,
		Password: users[0].Password,
	}).Expect().Status(http.StatusOK).JSON().Object()

	user1_token := r.Value("token").String().Raw()

	r = e.POST("/api/auth/login").WithJSON(model.LogInDTO{
		Email:    users[1].Email,
		Password: users[1].Password,
	}).Expect().Status(http.StatusOK).JSON().Object()

	user2_token := r.Value("token").String().Raw()

	r = e.POST("/api/auth/login").WithJSON(model.LogInDTO{
		Email:    users[2].Email,
		Password: users[2].Password,
	}).Expect().Status(http.StatusOK).JSON().Object()

	user3_token := r.Value("token").String().Raw()

	e.GET("/api/notes").Expect().Status(http.StatusUnauthorized)

	e.GET("/api/notes").WithHeader("Authorization", "Bearer <bad token>").Expect().Status(http.StatusUnauthorized)

	e.GET("/api/notes/").WithHeader("Authorization", "Bearer "+user1_token).Expect().Status(http.StatusOK)

	// Create note
	note1 := e.POST("/api/notes/").WithHeader("Authorization", "Bearer "+user1_token).WithJSON(model.NoteDTO{
		Title:   "title 1",
		Content: "description 1",
	}).Expect().Status(http.StatusOK).JSON().Object()

	note1_id := int(note1.Value("id").Number().Raw())

	e.GET("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user1_token).Expect().Status(http.StatusOK)

	// Share note to user2
	// It should not be visible to user2 before sharing
	e.GET("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user2_token).Expect().Status(http.StatusNotFound)

	e.POST("/api/notes/"+strconv.Itoa(note1_id)+"/share").WithHeader("Authorization", "Bearer "+user1_token).WithJSON(model.NoteShareDTO{
		NoteID:     int32(note1_id),
		SharedWith: users[1].Email,
	}).Expect().Status(http.StatusOK)

	// Now it should be visible
	e.GET("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user2_token).Expect().Status(http.StatusOK)

	// User3 still shouldn't have the access.
	e.GET("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user3_token).Expect().Status(http.StatusNotFound)

	// Search notes
	e.POST("/api/notes/").WithHeader("Authorization", "Bearer "+user1_token).WithJSON(model.NoteDTO{
		Title:   "title 2",
		Content: "description 2 with some long text",
	}).Expect().Status(http.StatusOK)

	e.POST("/api/notes/").WithHeader("Authorization", "Bearer "+user1_token).WithJSON(model.NoteDTO{
		Title:   "title 3",
		Content: "description 3 with some long text",
	}).Expect().Status(http.StatusOK)

	e.GET("/api/notes/search").WithQuery("q", "with text").WithHeader("Authorization", "Bearer "+user1_token).Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(2)

	// Delete note
	// Other user should not be able to delete the note even if it is shared with them.
	e.DELETE("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user2_token).Expect().Status(http.StatusForbidden)
	e.DELETE("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user3_token).Expect().Status(http.StatusForbidden)

	// Delete note by owner

	e.DELETE("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user1_token).Expect().Status(http.StatusOK)

	// It should not be visible anymore

	e.GET("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user1_token).Expect().Status(http.StatusNotFound)
	e.GET("/api/notes/"+strconv.Itoa(note1_id)).WithHeader("Authorization", "Bearer "+user2_token).Expect().Status(http.StatusNotFound)
}
