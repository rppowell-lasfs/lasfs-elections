package api

import (
	"election/storage"
	"election/types"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type GinHandler struct {
	storage storage.StorageInterface
}

func NewGinHandler(storage storage.StorageInterface) *GinHandler {
	return &GinHandler{storage: storage}
}

func (h *GinHandler) Signup(c *gin.Context) {
	type CreateLASFSMemberPayload struct {
		Name     string
		Password string
	}
	var body CreateLASFSMemberPayload
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error from bcrypt GenerateFromPassword",
		})
	}
	user := types.LASFSMember{Name: body.Name, Password: string(hash)}
	id, err := h.storage.CreateLASFSMember(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"name": id,
		})
	}
}

func (h *GinHandler) Login(c *gin.Context) {
	type LoginLASFSMemberPayload struct {
		Name     string
		Password string
	}
	var body LoginLASFSMemberPayload
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user := h.storage.GetLASFSMemberByName(body.Name)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid name or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid name or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Name,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	hmacSampleSecret := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error creating token",
		})
		return
	}
	use_cookie := true

	if use_cookie {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	} else {
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func (h *GinHandler) RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	hmacSampleSecret := os.Getenv("SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if time.Now().Unix() > claims["exp"].(int64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// fmt.Println(claims["sub"], claims["exp"])
		var member types.LASFSMember = *(h.storage.GetLASFSMemberByName(claims["id"].(string)))
		c.Set("user", member.Name)
		c.Next()
	} else {
		// fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

func (h *GinHandler) Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": "validated",
		"user":    user,
	})
}

func (h *GinHandler) GetLASFSElections(c *gin.Context) {
	elections := h.storage.GetLASFSElectionsIDs()
	c.JSON(http.StatusOK, gin.H{
		"elections": elections,
	})
}

func RunGIN() {
	s := storage.NewDevStorage()
	ginHandler := NewGinHandler(s)

	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/signup", ginHandler.Signup)
	r.POST("/login", ginHandler.Login)
	// r.GET("/validate", RequireAuth, Validate)
	r.GET("/elections", ginHandler.GetLASFSElections)
	r.Run()
}
