package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/malcolmmaima/maimabank/token"
)

// Constants for authorization header keys and types.
const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a middleware function that performs token authentication.
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get authorization header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "authorization header not provided"})
			return
		}

		// Split authorization header into fields (type and token)
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization header format"})
			return
		}

		// Extract the authorization type and check if it is "bearer"
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unsupported authorization type"})
			return
		}

		// Extract the access token from the authorization header
		accessToken := fields[1]

		// Verify the access token using the token maker
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		// Set the authorization payload in the context for further use
		ctx.Set(authorizationPayloadKey, payload)

		// Call the next middleware or handler in the chain
		ctx.Next()
	}
}
