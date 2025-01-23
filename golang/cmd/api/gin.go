package api

import (
	"election/storage"
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
	var body CreateLASFSMemberRequestPayload
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "Failed to read body",
		})
		return
	}
	// TODO CreateLASFSMemberPayload JSON validator here

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponsePayload{
			ErrorMessage: "Error from bcrypt GenerateFromPassword",
		})
	}
	user, err := h.storage.CreateLASFSMember(body.Name, string(hash))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: err.Error(),
		})
	} else {
		c.JSON(http.StatusCreated, CreateLASFSMemberResponsePayload{
			ID: user.ID,
		})
	}
}

func (h *GinHandler) Login(c *gin.Context) {
	var body LASFSMemberLoginRequestPayload
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "Failed to read body",
		})
		return
	}
	// TODO LoginLASFSMemberPayload JSON validator here

	user, err := h.storage.GetLASFSMemberByName(body.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "Invalid name or password",
		})
		return
	}

	bcrypterr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if bcrypterr != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "Invalid name or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Name,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	//  TODO: fix config.Envs godotenv to handle paths for testing
	//	hmacSampleSecret := config.Envs.APIConfig.BcryptSecret
	hmacSampleSecret := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponsePayload{
			ErrorMessage: "error creating token",
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

	//  TODO: fix config.Envs godotenv to handle paths for testing
	//	hmacSampleSecret := config.Envs.APIConfig.BcryptSecret
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
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if time.Now().Unix() > claims["exp"].(int64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// fmt.Println(claims["sub"], claims["exp"])

		member, err := h.storage.GetLASFSMemberByName(claims["id"].(string))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

		}
		// var member types.LASFSMember = *(h.storage.GetLASFSMemberByName(claims["id"].(string)))
		c.Set("user", member.Name)
		c.Next()
	} else {
		// fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

func (h *GinHandler) Validate(c *gin.Context) {
	// TODO: Work In Progress here - Validate Login
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": "validated",
		"user":    user,
	})
}

func (h *GinHandler) GetLASFSElections(c *gin.Context) {
	elections := h.storage.GetLASFSElectionsIDs()
	c.JSON(http.StatusOK, GetLASFSElectionsResponsePayload{
		LASFSElections: elections,
	})
}

func (h *GinHandler) GetLASFSElection(c *gin.Context) {
	election_id := c.Param("election_id")
	election, err := h.storage.GetLASFSElectionByID(election_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "election not found",
		})
		return
	} else {
		nominees := election.Nominees
		if election.Nominees == nil {
			nominees = make([]string, 0)
		}
		c.JSON(http.StatusOK, gin.H{
			"election:": GetLASFSElectionResponsePayload{
				ElectionID:       election.ID,
				ElectionPosition: election.Position,
				ElectionStatus:   election.Status,
				Nominees:         nominees,
			},
		})
	}
}

func (h *GinHandler) GetBallotForElectionByMember(c *gin.Context) {
	election_id := c.Param("election_id")
	member_id := c.Param("member_id")

	election, err := h.storage.GetLASFSElectionByID(election_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "election not found",
		})
		return
	} else {
		ballot, err := election.GetLASFSMemberBallot(member_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponsePayload{
				ErrorMessage: "ballot not found",
			})
		} else {
			c.JSON(http.StatusOK, GetLASFSBallotResponsePayload{
				Election: election.ID,
				VoterID:  ballot.VoterName,
				Nominees: ballot.NomineeVotes(),
			},
			)
		}
	}
}

func (h *GinHandler) PostBallotForElectionByMember(c *gin.Context) {
	election_id := c.Param("election_id")
	member_id := c.Param("member_id")
	var body []string
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "Failed to read body",
		})
		return
	}
	ballot, err := h.storage.NewLASFSBallot(member_id, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponsePayload{
			ErrorMessage: "error making ballot",
		})
		return
	}

	_, err = h.storage.AddLASFSBallot(election_id, *ballot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponsePayload{
			ErrorMessage: "error adding ballot",
		})
	}
	c.JSON(http.StatusAccepted, GetLASFSBallotResponsePayload{
		Election: election_id,
		VoterID:  ballot.VoterName,
		Nominees: ballot.NomineeVotes(),
	},
	)
}

func SetupGin(ginHandler *GinHandler, r *gin.Engine) {
	r.POST("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/signup", ginHandler.Signup)
	r.POST("/login", ginHandler.Login)
	// r.GET("/validate", RequireAuth, Validate)
	r.GET("/elections", ginHandler.GetLASFSElections)
	r.GET("/election/:election_id", ginHandler.GetLASFSElection)
	// r.GET("/votes/:election_id/", ginHandler.GetVotesForElection)
	r.GET("/vote/:election_id/:member_id", ginHandler.GetBallotForElectionByMember)
	r.POST("/vote/:election_id/:member_id", ginHandler.PostBallotForElectionByMember)
}

func RunGIN() {
	s := storage.NewDevStorage()
	r := gin.Default()
	ginHandler := NewGinHandler(s)
	SetupGin(ginHandler, r)
	r.Run()
}
