package helper

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sahil/jwt-auth-go/configs"
	"github.com/sahil/jwt-auth-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetails
type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = configs.OpenCollection(configs.Client, "user")

var SECRET_KEY string = os.Getenv(configs.SecretKey())

// GenerateAllTokens generates both teh detailed token and refresh token
func GenerateAllTokens(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

// ValidateToken validates the jwt token
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = err.Error()
		return
	}

	return claims, msg
}

// // UpdateAllTokens renews the user tokens when they login
// func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) (models.User, error) {

// 	var updateObj primitive.D

// 	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
// 	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
// 	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
// 	Updated_at := now.Format("2006-01-02 15:04:05")
// 	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

// 	upsert := true
// 	filter := bson.M{"user_id": userId}
// 	opt := options.FindOneAndUpdateOptions{
// 		Upsert: &upsert,
// 	}
// 	opt.SetReturnDocument(options.After)
// 	var updatedUser models.User
// 	err := userCollection.FindOneAndUpdate(
// 		context.TODO(),
// 		filter,
// 		bson.D{{Key: "$set", Value: updateObj}},
// 		&opt,
// 	).Decode(&updatedUser)

// 	return updatedUser, err
// }

// Update user
func UpdateUserByID(token string, refresh string, id string) (models.User, error) {
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refresh})

	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	Updated_at := now.Format("2006-01-02 15:04:05")
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})
	upsert := true
	filter := bson.M{"user_id": id}
	opt := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}
	opt.SetReturnDocument(options.After)
	var updatedUser models.User

	err := userCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	).Decode(&updatedUser)

	return updatedUser, err
}
