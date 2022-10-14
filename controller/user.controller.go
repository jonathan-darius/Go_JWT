package controller

import (
	_ "Latihan_Mongo/httputil"
	"Latihan_Mongo/model"
	"Latihan_Mongo/service"
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
	"net/http"
)

type UserController struct {
	UserService service.UserService
}

func New(userservice service.UserService) UserController {
	return UserController{
		UserService: userservice,
	}
}

var imageTmp model.Image

// CreateUser godoc
// @Summary      Add User
// @Description  post User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body	model.User true "Add User"
// @Success      200  {string}	{"message": "success"}
// @Failure      400  {object}  httputil.HTTPError
// @Router       /user/create [post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	var in model.User_IN
	err := ctx.Bind(&user)
	in.Name = user.Name
	in.Age = user.Age
	in.Address = user.Address
	in.Path = user.Path
	if err != nil {
		ctx.String(500, "error: %v\n", err)
		return
	}
	err = uc.UserService.CreateUser(&in)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// GetUser godoc
// @Summary      Show an User
// @Description  get string by name
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        name   path      string  true  "User Name"
// @Success      200  {object}  model.User
// @Failure      400  {object}  httputil.HTTPError
// @Router       /user/get/{name} [get]
func (uc *UserController) GetUser(ctx *gin.Context) {
	name := ctx.Param("name")
	user, err := uc.UserService.GetUser(&name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// GetAll godoc
// @Summary      List User
// @Description  get all user data
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.User
// @Failure      400  {object}  httputil.HTTPError
// @Router       /user/getall [get]
func (uc *UserController) GetAll(ctx *gin.Context) {
	allUser, err := uc.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, allUser)
}

// UpdateUser godoc
// @Summary      Update	 User
// @Description  Update by json User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      model.User  true  "Update account"
// @Success      200      {string}	{"message": "success"}
// @Failure      400      {object}  httputil.HTTPError
// @Router       /user/update [patch]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// DeleteUser godoc
// @Summary      Delete an User
// @Description  Delete User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        name   	path      	string  true  "delete name"
// @Success      200      	{string}	{"message": "success"}
// @Failure      400      	{object}  	httputil.HTTPError
// @Router       /user/delete/{name} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	usermane := ctx.Param("name")
	err := uc.UserService.DeleteUser(&usermane)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (uc *UserController) UploadImage(ctx *gin.Context) {
	gambar, err := imageupload.Process(ctx.Request, "file")

	imageTmp.ImageFile = gambar
	if err != nil {
		panic(err)
	}
	//buat fungsi ke user service
	imageTmp.ImageFile.Write(ctx.Writer)
}

func (uc *UserController) ViewImage(ctx *gin.Context) {
	if imageTmp.ImageFile == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	imageTmp.ImageFile.Write(ctx.Writer)
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userRoute := rg.Group("/user")
	userRoute.POST("/create", uc.CreateUser)
	userRoute.GET("/get/:name", uc.GetUser)
	userRoute.GET("/getall", uc.GetAll)
	userRoute.PATCH("/update", uc.UpdateUser)
	userRoute.DELETE("/delete/:name", uc.DeleteUser)
	userRoute.POST("/upload", uc.UploadImage)
	userRoute.GET("/image_view", uc.ViewImage)
}
