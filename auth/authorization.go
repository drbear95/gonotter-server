package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/drbear95/gonotter-server/model"
	"github.com/drbear95/gonotter-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userid string) (*TokenDetails, error) {
	c := utils.Config{}
	c.GetConfig()
	secret := c.JWTSecret
	refreshSecret := c.JWTRefreshSecret

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(userid string, td *TokenDetails) error {
	refreshTokenDetails := model.NewRefreshTokenDetails(userid, td.RefreshUuid, td.RtExpires)
	accessTokenDetails := model.NewAccessTokenDetails(userid, td.AccessUuid, td.AtExpires)

	var errAccess = mgm.Coll(accessTokenDetails).Create(refreshTokenDetails)
	if errAccess != nil {
		return errAccess
	}

	var errRefresh = mgm.Coll(refreshTokenDetails).Create(refreshTokenDetails)
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func SignIn(c *gin.Context) {
	var u model.User
	coll := mgm.Coll(&u)
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	result := model.User{}
	err := coll.First(bson.M{"name": u.Name, "$and": []bson.M{{"password": u.Password}}}, &result)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	ts, err := CreateToken(result.ID.Hex())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := CreateAuth(result.ID.Hex(), ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func SignUp(c *gin.Context) {
	var u model.User
	coll := mgm.Coll(&u)

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	var err error

	err = validateUser(u)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprint(err))
		return
	}

	err = coll.Create(&u)

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	} else {
		c.JSON(http.StatusOK, "User has been created")
	}
}

func validateUser(user model.User) error {
	var err error

	if user.Name == ""{
		return errors.New("username is empty")
	}

	if user.Password == ""{
		return errors.New("password is empty")
	}

	if user.Email == ""{
		return errors.New("email is empty")
	}

	var result []model.User

	err = mgm.Coll(&model.User{}).SimpleFind(&result, bson.M{"name": bson.M{operator.Eq: user.Name}})

	if err != nil {
		return nil
	} else {
		if result == nil {
			return nil
		} else {
			return errors.New("username already exists")
		}
	}
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	c := utils.Config{}
	c.GetConfig()
	secret := c.JWTSecret

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*model.AccessTokenDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, err
		}
		return &model.AccessTokenDetails{
			UserID:    userId,
			Uuid:      accessUuid,
			ExpiresAt: int64(exp),
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *model.AccessTokenDetails) (string, error) {
	var atd model.AccessTokenDetails
	coll := mgm.Coll(&atd)

	err := coll.First(bson.M{"user_id": atd.UserID}, &atd)

	if err != nil {
		return "", err
	}

	return atd.UserID, nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserId(context context.Context) (*primitive.ObjectID, error) {
	authDetails, ok := context.Value("auth_details").(*model.AccessTokenDetails)

	if !ok {
		return nil, errors.New("auth details not ok")
	}

	userId, err := primitive.ObjectIDFromHex(authDetails.UserID)

	if !ok {
		return nil, err
	}

	return &userId, nil
}
