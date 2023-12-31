package controller

import (
	"net/http"
	"os"
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
	GetTopBooks(c *gin.Context)
	GetImage(c *gin.Context)
}

type bookController struct {
	bookService services.BookService
	jwtService  services.JWTService
	userService services.UserService
}

func NewBookController(bs services.BookService, jwt services.JWTService, us services.UserService) BookController {
	return &bookController{
		bookService: bs,
		jwtService:  jwt,
		userService: us,
	}
}

func (bc *bookController) CreateBook(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := bc.jwtService.GetIDByToken(token)
	if err != nil {
		res := utils.BuildResponseFailed("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	user, err := bc.userService.GetUserByID(ctx, userID)
	if err != nil {
		res := utils.BuildResponseFailed("Gagal Memproses Request", "userID tidak valid", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if user.Role != "admin" {
		res := utils.BuildResponseFailed("Tidak memiliki Akses", "Role Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	var req dto.BookCreateRequest
	req.UserID = userID
	req.Title = ctx.PostForm("title")
	if checkTitle, _ := bc.bookService.CheckTitle(ctx.Request.Context(), req.Title); checkTitle {
		res := utils.BuildResponseFailed("Judul Sudah Terdaftar", "failed", utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req.Desc = ctx.PostForm("description")
	if req.Desc == "" || req.Title == "" {
		res := utils.BuildResponseFailed("Failed to retrieve title/desc", "", utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	thumbnail, err := ctx.FormFile("thumbnail")
	if err != nil {
		res := utils.BuildResponseFailed("Failed to retrieve thumbnail", err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req.Thumbnail = thumbnail

	for i := 1; ; i++ {
		filesUploaded := false
		index := 1
		for j := 1; ; j++ {

			photo, err := ctx.FormFile("Page[" + strconv.Itoa(i) + "][" + strconv.Itoa(j) + "]")
			if err != nil {
				break
			}

			if photo == nil {
				break
			}

			var medias dto.MediaRequest
			medias.Media = photo
			medias.Index = index
			medias.Page = i
			req.MediaRequest = append(req.MediaRequest, medias)

			filesUploaded = true
			index++
		}

		if !filesUploaded {
			index--
			break
		}
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

func (bc *bookController) GetTopBooks(c *gin.Context) {
	result, err := bc.bookService.GetTopBooks(c.Request.Context())

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
	page := ctx.Query("page")
	if page == "" {
		page = "1"
	}
	Books, err := bc.bookService.GetBookPages(ctx, id, page)
	if err != nil {
		res := utils.BuildResponseFailed("Gagal Mendapatkan Detail Buku", err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponseSuccess("Berhasil Mendapatkan Project", Books)
	ctx.JSON(http.StatusOK, res)
}

func (bc *bookController) GetImage(ctx *gin.Context) {
	path := ctx.Param("path")
	dirname := ctx.Param("dirname")
	filename := ctx.Param("filename")

	imagePath := "storage/" + path + "/" + dirname + "/" + filename

	_, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		ctx.JSON(400, gin.H{
			"message": "image not found",
		})
		return
	}

	ctx.File(imagePath)
	
}
