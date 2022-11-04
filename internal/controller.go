package internal

import (
	"html"
	"net/http"
	"skill-test-dans/shared"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Controllers interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetJobList(c *gin.Context)
	GetJobDetail(c *gin.Context)
}

type controllers struct {
	db   *gorm.DB
	data []*Job
}

func NewController(db *gorm.DB, data []*Job) Controllers {
	return &controllers{
		db:   db,
		data: data,
	}
}

func (s *controllers) Register(c *gin.Context) {
	var input InputRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := shared.GeneratePassword(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{
		Password: hashedPassword,
		Username: html.EscapeString(strings.TrimSpace(input.Username)),
	}

	err = s.db.Create(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func (s *controllers) Login(c *gin.Context) {
	var input InputRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}
	err := s.db.Model(User{}).Where("username = ?", input.Username).Take(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	err = shared.VerifyPassword(input.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	token, err := GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token."})
		return
	}
	whiteListTokens = append(whiteListTokens, string(token))

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *controllers) Logout(c *gin.Context) {
	token := ExtractToken(c)
	whiteListTokens = shared.RemoveString(whiteListTokens, token)

	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

func (s *controllers) GetJobList(c *gin.Context) {
	user_id, err := ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u User
	if err := s.db.First(&u, user_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.Password = "********"

	var paginate ParamPaginate
	if err = c.ShouldBindQuery(&paginate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if paginate.Page != nil {
		page = *paginate.Page
	}
	if paginate.Limit != nil {
		limit = *paginate.Limit
	}

	result := Paginate(s.data, page, limit)
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": result})
}

func (s *controllers) GetJobDetail(c *gin.Context) {
	param := c.Params.ByName("id")
	user_id, err := ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u User
	if err := s.db.First(&u, user_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.Password = "********"

	data := s.data
	jobDetail := map[string]*Job{}
	for _, v := range data {
		jobDetail[v.ID] = v
	}

	var result *Job
	if v, ok := jobDetail[param]; ok {
		result = v
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": result})
}
