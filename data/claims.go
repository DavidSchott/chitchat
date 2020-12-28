package data

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	expirationConstantMinutes                   int = 60  // TODO: Change to 10 minutes
	minimumElapsedTimeToRefreshToleranceMinutes int = 300 // TODO: Change to 5 minutes
)

// Claims is a model that represents JSON web tokens used for authentication by users
type Claims struct {
	Username string `json:"username"`
	RoomID   string `json:"room_id,omitempty"`
	jwt.StandardClaims
}

//EncodeJWT will generate a jwt token based
func EncodeJWT(c *ChatEvent, cr *ChatRoom, secretKey string) (tokenString string, err error) {
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(time.Duration(expirationConstantMinutes) * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: c.User,
		RoomID:   cr.ID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(secretKey)
	// Create the JWT string
	tokenString, err = token.SignedString(jwtKey)
	return
}

// ParseJWT parses a JWT and stores Claims object in c
func ParseJWT(tokenString string, c *Claims, secretKey string) (err error) {
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	switch err {
	// TODO: What are some other useful cases?
	case jwt.ErrSignatureInvalid:
		err = &APIError{
			Code:  401,
			Field: "signature",
		}
	case nil:
		// Check signing algorithm is as expected:
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return &APIError{
				Code:  402,
				Field: "signing method",
			}
		}
		if !tkn.Valid {
			return &APIError{
				Code:  403,
				Field: "token",
			}
		}
	default:
		err = &APIError{
			Code:  403,
			Field: "token",
		}
	}

	return
}

//RefreshJWT will refresh return a signed token with new expiration time
func (c Claims) RefreshJWT(secretKey string) (tokenString string, err error) {
	// Ensure enough time has elapsed since last token was generated
	if time.Until(time.Unix(c.ExpiresAt, 0)) > time.Duration(minimumElapsedTimeToRefreshToleranceMinutes)*time.Minute {
		return "", &APIError{
			Code:  403,
			Field: "token",
		}
	}
	// Now, create a new token with a renewed expiration time
	newExpirationTime := time.Now().Add(time.Duration(expirationConstantMinutes) * time.Minute)
	c.ExpiresAt = newExpirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secretKey))
}
