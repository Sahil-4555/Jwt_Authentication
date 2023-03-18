package controllers

import (
	"context"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/sahil/jwt-auth-go/configs"
	helper "github.com/sahil/jwt-auth-go/helpers"
	"github.com/sahil/jwt-auth-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.OpenCollection(configs.Client, "jwt-user")
var validate = validator.New()

// Insert user into DB
func InsertUser(u *models.User) (interface{}, error) {
	result, err := userCollection.InsertOne(context.TODO(), u)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

// FindUserCountByMail to verify
func FindUserCountByMail(email string) (int64, error) {
	count, err := userCollection.CountDocuments(context.TODO(), bson.M{"email": email})
	if err != nil {
		return count, err
	}
	return count, nil
}

// Find User By Email Search
func FindUserByEmail(ctx context.Context, email *string) (models.User, error) {
	var foundUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return foundUser, err
	}
	return foundUser, err
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

// CreateUser is the api used to tget a single user
func (con *Controllers) SignUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	count, err := FindUserCountByMail(*user.Email)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		return
	}

	password := HashPassword(*user.Password)
	user.Password = &password
	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	user.Created_at = now.Format("2006-01-02 15:04:05")
	user.Updated_at = now.Format("2006-01-02 15:04:05")
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken
	_, insertErr := InsertUser(&user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user item was not created"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})

}

// Login is the api used to tget a single user
func (c *Controllers) Login(ctx *gin.Context) {
	ctxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"data": "Inavlid JSON Provided"})
	}

	foundUser, err := FindUserByEmail(ctxt, user.Email)
	defer cancel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data": "Error Occured While Checking For The Email"})
		return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if !passwordIsValid {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)
	updateUser, err := helper.UpdateUserByID(token, refreshToken, foundUser.User_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data": "error occured while logging in",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": updateUser,
	})
}
