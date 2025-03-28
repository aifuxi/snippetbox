package main

import (
	"errors"
	"fmt"
	"github.com/aifuxi/snippetbox/internal/models"
	"github.com/aifuxi/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// 获取所有params
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int

	validator.Validator
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// 1. 解析表单提交的数据，失败返回400错误
	err := r.ParseForm() // 会把解析好的数据存在 r.PostForm 这个map中
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// 2. 从 r.PostForm这个map中提取相应的数据，并保存到数据库中
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.serverError(w, err)
		return
	}

	// 3. 校验参数
	form := snippetCreateForm{
		Title:   title,
		Content: content,
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(title, 100), "title", "This field cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(content), "content", "This field cannot be blank")
	permittedValues := []int{1, 7, 365}

	var strValues []string
	for _, value := range permittedValues {
		strValues = append(strValues, strconv.Itoa(value))
	}
	result := strings.Join(strValues, ", ")

	form.CheckField(
		validator.PermittedInt(expires, permittedValues...),
		"content",
		fmt.Sprintf("This field must equal %s", result),
	)

	// 如果有错，直接返回
	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 重定向
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
