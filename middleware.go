package community_authentication_middleware

import (
	"context"
	"github.com/satori/go.uuid"
	cd "github.com/214alphadev/community-bl"
	"net/http"
)

type AuthenticateMemberMiddleware = func(next http.Handler) http.Handler

const middlewareAuthKey = "community-member"

func NewAuthenticateMemberMiddleware(community cd.CommunityInterface) AuthenticateMemberMiddleware {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// member
			member, err := community.GetMemberByAccessToken(r.Header.Get("Authorization-Bearer"))

			switch err {
			case nil:
				emptyUUID := uuid.UUID{}.String()
				if member.ID.String() == emptyUUID {
					next.ServeHTTP(w, r)
					return
				}
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), middlewareAuthKey, member.ID)))
			default:
				next.ServeHTTP(w, r)
			}

		})

	}

}

func GetAuthenticateMember(ctx context.Context) *cd.MemberIdentifier {
	switch v := ctx.Value(middlewareAuthKey).(type) {
	case uuid.UUID:
		return &v
	case *uuid.UUID:
		return v
	default:
		return nil
	}
}
