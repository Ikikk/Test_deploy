package controller

import (
	"net/http"
	"strconv"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/gin-gonic/gin"
)

type BookController interface {
	CreateBook(c *gin.Context)
	GetAllBooks(c *gin.Context)
	GetBookPages(c *gin.Context)
}

type bookController struct {
	bookService services.BookService
}

func NewBookController(bs services.BookService) BookController {
	return &bookController{
		bookService: bs,
	}
}

func (bc *bookController) CreateBook(ctx *gin.Context) {
	var req dto.BookCreateRequest

	ctx.PostFormArray("pages")
	req.Title = ctx.Request.PostForm.Get("title")

	thumbnail, err := ctx.FormFile("thumbnail")
	if err != nil {
		res := utils.BuildResponseFailed("Failed to save thumbnail", err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	req.Thumbnail = thumbnail

	for i := 0; ; i++ {
		photo, err := ctx.FormFile("Page[" + strconv.Itoa(i) + "][pages]")
		if err != nil {
			break
		}

		if photo == nil {
			break
		}

		var Pages dto.PagesRequest
		Pages.Pages = photo
		req.PagesRequest = append(req.PagesRequest, Pages)
	}

	Page, err := bc.bookService.CreateBook(ctx, req)
	if err != nil {
		res := utils.BuildResponseFailed("Failed to create Books", err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess("Successfully create Books", Page)
	ctx.JSON(http.StatusOK, res)

}

func (bc *bookController) GetAllBooks(c *gin.Context) {
	result, err := bc.bookService.GetAllBooks(c.Request.Context())

	if err != nil {
		res := utils.BuildResponseFailed("Gagal Mendapatkan List Buku", err.Error(), utils.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess("Berhasil Mendapatkan List Buku", result)
	c.JSON(http.StatusOK, res)
}

func (bc *bookController) GetBookPages(ctx *gin.Context) {
	id := ctx.Param("book_id")

	Books, err := bc.bookService.GetBookPages(ctx, id)
	if err != nil {
		res := utils.BuildResponseFailed("Gagal Mendapatkan Detail Buku", err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess("Berhasil Mendapatkan Project", Books)
	ctx.JSON(http.StatusOK, res)

}
