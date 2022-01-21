package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type (
	mockAuthService struct {
		authService

		store         store.Storer
		settings      *types.AppSettings
		notifications service.AuthNotificationService

		providerValidator func(string) error
	}

	mockNotificationService struct {
		settings *types.AppSettings
		opt      options.AuthOpt
	}

	mockSession struct {
		Options *sessions.Options
		get     func(r *http.Request, name string) (*sessions.Session, error)
		new     func(r *http.Request, name string) (*sessions.Session, error)
		save    func(r *http.Request, w http.ResponseWriter, s *sessions.Session) error
	}

	oauth2ServiceMocked struct {
		getRedirectURI             func(req *server.AuthorizeRequest, data map[string]interface{}) (string, error)
		checkResponseType          func(rt oauth2.ResponseType) bool
		checkCodeChallengeMethod   func(ccm oauth2.CodeChallengeMethod) bool
		validationAuthorizeRequest func(r *http.Request) (*server.AuthorizeRequest, error)
		getAuthorizeToken          func(ctx context.Context, req *server.AuthorizeRequest) (oauth2.TokenInfo, error)
		getAuthorizeData           func(rt oauth2.ResponseType, ti oauth2.TokenInfo) map[string]interface{}
		handleAuthorizeRequest     func(w http.ResponseWriter, r *http.Request) error
		validationTokenRequest     func(r *http.Request) (oauth2.GrantType, *oauth2.TokenGenerateRequest, error)
		checkGrantType             func(gt oauth2.GrantType) bool
		getAccessToken             func(ctx context.Context, gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (oauth2.TokenInfo, error)
		getTokenData               func(ti oauth2.TokenInfo) map[string]interface{}
		handleTokenRequest         func(w http.ResponseWriter, r *http.Request) error
		getErrorData               func(err error) (map[string]interface{}, int, http.Header)
		bearerAuth                 func(r *http.Request) (string, bool)
		validationBearerToken      func(r *http.Request) (oauth2.TokenInfo, error)
	}

	testingExpect struct {
		name        string
		payload     interface{}
		link        string
		err         string
		template    string
		alerts      []request.Alert
		userService userService
		fn          func(*settings.Settings)
	}

	userServiceMocked struct {
		update    func(context.Context, *types.User) (*types.User, error)
		findByAny func(context.Context, interface{}) (*types.User, error)
	}

	authServiceMocked struct {
		external                          func(context.Context, types.ExternalAuthUser) (u *types.User, err error)
		internalSignUp                    func(context.Context, *types.User, string) (u *types.User, err error)
		internalLogin                     func(context.Context, string, string) (u *types.User, err error)
		setPassword                       func(context.Context, uint64, string) (err error)
		changePassword                    func(context.Context, uint64, string, string) (err error)
		createPassword                    func(context.Context, uint64, string) (err error)
		validateEmailConfirmationToken    func(context.Context, string) (user *types.User, err error)
		validatePasswordResetToken        func(context.Context, string) (user *types.User, err error)
		validatePasswordCreateToken       func(context.Context, string) (user *types.User, err error)
		sendEmailAddressConfirmationToken func(context.Context, *types.User) (err error)
		sendPasswordResetToken            func(context.Context, string) (err error)
		passwordSet                       func(context.Context, string) bool
		getProviders                      func() types.ExternalAuthProviderSet
		validateTOTP                      func(context.Context, string) (err error)
		configureTOTP                     func(context.Context, string, string) (u *types.User, err error)
		removeTOTP                        func(context.Context, uint64, string) (u *types.User, err error)
		sendEmailOTP                      func(context.Context) (err error)
		configureEmailOTP                 func(context.Context, uint64, bool) (u *types.User, err error)
		validateEmailOTP                  func(context.Context, string) (err error)
	}
)

//
// Mocking userService
//
func (u userServiceMocked) Update(ctx context.Context, user *types.User) (*types.User, error) {
	return u.update(ctx, user)
}

func (u userServiceMocked) FindByAny(ctx context.Context, any interface{}) (*types.User, error) {
	return u.findByAny(ctx, any)
}

//
// Mocking authService
//
func (s authServiceMocked) External(ctx context.Context, profile types.ExternalAuthUser) (u *types.User, err error) {
	return s.external(ctx, profile)
}

func (s authServiceMocked) InternalSignUp(ctx context.Context, input *types.User, password string) (u *types.User, err error) {
	return s.internalSignUp(ctx, input, password)
}

func (s authServiceMocked) InternalLogin(ctx context.Context, email string, password string) (u *types.User, err error) {
	return s.internalLogin(ctx, email, password)
}

func (s authServiceMocked) SetPassword(ctx context.Context, userID uint64, password string) (err error) {
	return s.setPassword(ctx, userID, password)
}

func (s authServiceMocked) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (err error) {
	return s.changePassword(ctx, userID, oldPassword, newPassword)
}

func (s authServiceMocked) ValidateEmailConfirmationToken(ctx context.Context, token string) (user *types.User, err error) {
	return s.validateEmailConfirmationToken(ctx, token)
}

func (s authServiceMocked) ValidatePasswordResetToken(ctx context.Context, token string) (user *types.User, err error) {
	return s.validatePasswordResetToken(ctx, token)
}

func (s authServiceMocked) ValidatePasswordCreateToken(ctx context.Context, token string) (user *types.User, err error) {
	return s.validatePasswordCreateToken(ctx, token)
}

func (s authServiceMocked) PasswordSet(ctx context.Context, email string) (is bool) {
	return s.passwordSet(ctx, email)
}

func (s authServiceMocked) SendEmailAddressConfirmationToken(ctx context.Context, u *types.User) (err error) {
	return s.sendEmailAddressConfirmationToken(ctx, u)
}

func (s authServiceMocked) SendPasswordResetToken(ctx context.Context, email string) (err error) {
	return s.sendPasswordResetToken(ctx, email)
}

func (s authServiceMocked) GetProviders() types.ExternalAuthProviderSet {
	return s.getProviders()
}

func (s authServiceMocked) ValidateTOTP(ctx context.Context, code string) (err error) {
	return s.validateTOTP(ctx, code)
}

func (s authServiceMocked) ConfigureTOTP(ctx context.Context, secret string, code string) (u *types.User, err error) {
	return s.configureTOTP(ctx, secret, code)
}

func (s authServiceMocked) RemoveTOTP(ctx context.Context, userID uint64, code string) (u *types.User, err error) {
	return s.removeTOTP(ctx, userID, code)
}

func (s authServiceMocked) SendEmailOTP(ctx context.Context) (err error) {
	return s.sendEmailOTP(ctx)
}

func (s authServiceMocked) ConfigureEmailOTP(ctx context.Context, userID uint64, enable bool) (u *types.User, err error) {
	return s.configureEmailOTP(ctx, userID, enable)
}

func (s authServiceMocked) ValidateEmailOTP(ctx context.Context, code string) (err error) {
	return s.validateEmailOTP(ctx, code)
}

func (s authServiceMocked) LoadRoleMemberships(ctx context.Context, u *types.User) error {
	// no-op for now
	return nil
}

//
// Mocking oauth2Service
//
func (s *oauth2ServiceMocked) GetRedirectURI(req *server.AuthorizeRequest, data map[string]interface{}) (string, error) {
	return s.GetRedirectURI(req, data)
}

func (s *oauth2ServiceMocked) CheckResponseType(rt oauth2.ResponseType) bool {
	return s.checkResponseType(rt)
}

func (s *oauth2ServiceMocked) CheckCodeChallengeMethod(ccm oauth2.CodeChallengeMethod) bool {
	return s.CheckCodeChallengeMethod(ccm)
}

func (s *oauth2ServiceMocked) ValidationAuthorizeRequest(r *http.Request) (*server.AuthorizeRequest, error) {
	return s.validationAuthorizeRequest(r)
}

func (s *oauth2ServiceMocked) GetAuthorizeToken(ctx context.Context, req *server.AuthorizeRequest) (oauth2.TokenInfo, error) {
	return s.getAuthorizeToken(ctx, req)
}

func (s *oauth2ServiceMocked) GetAuthorizeData(rt oauth2.ResponseType, ti oauth2.TokenInfo) map[string]interface{} {
	return s.GetAuthorizeData(rt, ti)
}

func (s *oauth2ServiceMocked) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return s.handleAuthorizeRequest(w, r)
}

func (s *oauth2ServiceMocked) ValidationTokenRequest(r *http.Request) (oauth2.GrantType, *oauth2.TokenGenerateRequest, error) {
	return s.validationTokenRequest(r)
}

func (s *oauth2ServiceMocked) CheckGrantType(gt oauth2.GrantType) bool {
	return s.checkGrantType(gt)
}

func (s *oauth2ServiceMocked) GetAccessToken(ctx context.Context, gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (oauth2.TokenInfo, error) {
	return s.getAccessToken(ctx, gt, tgr)
}

func (s *oauth2ServiceMocked) GetTokenData(ti oauth2.TokenInfo) map[string]interface{} {
	return s.getTokenData(ti)
}

func (s *oauth2ServiceMocked) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.handleTokenRequest(w, r)
}

func (s *oauth2ServiceMocked) GetErrorData(err error) (map[string]interface{}, int, http.Header) {
	return s.getErrorData(err)
}

func (s *oauth2ServiceMocked) BearerAuth(r *http.Request) (string, bool) {
	return s.bearerAuth(r)
}

func (s *oauth2ServiceMocked) ValidationBearerToken(r *http.Request) (oauth2.TokenInfo, error) {
	return s.validationBearerToken(r)
}

//
// Mocking authService
//
func (ma mockAuthService) ValidatePasswordResetToken(ctx context.Context, token string) (*types.User, error) {
	return &types.User{ID: 123}, nil
}

//
// Mocking authService
//
func (ma mockAuthService) ValidatePasswordCreateToken(ctx context.Context, token string) (*types.User, error) {
	return &types.User{ID: 123}, nil
}

func (ma mockAuthService) SendEmailOTP(ctx context.Context) error {
	return nil
}

func (ma mockAuthService) ValidateTOTP(ctx context.Context, code string) (err error) {
	err = nil
	return
}

//
// Mocking notification service
//
func (m mockNotificationService) EmailConfirmation(ctx context.Context, emailAddress string, token string) error {
	return nil
}

func (m mockNotificationService) PasswordReset(ctx context.Context, emailAddress string, token string) error {
	return nil
}

func (m mockNotificationService) PasswordCreate(token string) (string, error) {
	return "", nil
}

func (m mockNotificationService) EmailOTP(ctx context.Context, emailAddress string, code string) error {
	return nil
}

//
// Mocking gorilla session
//
func (ms mockSession) Get(r *http.Request, name string) (*sessions.Session, error) {
	return ms.get(r, name)
}
func (ms mockSession) New(r *http.Request, name string) (*sessions.Session, error) {
	return ms.new(r, name)
}
func (ms mockSession) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return ms.save(r, w, s)
}

//
// Helpers
//
func makeMockAuthService() *mockAuthService {
	service.DefaultAuthNotification = mockNotificationService{
		settings: service.CurrentSettings,
		opt:      options.AuthOpt{},
	}

	serviceAuth := service.Auth(service.AuthOptions{})

	svc := mockAuthService{
		authService: serviceAuth,
		settings:    service.CurrentSettings,
		providerValidator: func(s string) error {
			return nil
		},
	}

	return &svc
}

func makeMockUser() *types.User {
	u := &types.User{ID: 1, Username: "mock.user", Email: "mockuser@example.tld", Meta: &types.UserMeta{}}
	u.Meta.SecurityPolicy.MFA.EnforcedEmailOTP = true
	u.Meta.SecurityPolicy.MFA.EnforcedTOTP = false

	return u
}

func prepareClientAuthReq(h *AuthHandlers, req *http.Request, user *types.User) *request.AuthReq {
	s := &settings.Settings{}
	s.MultiFactor.EmailOTP.Enabled = true
	s.MultiFactor.EmailOTP.Enforced = true
	s.MultiFactor.TOTP.Enabled = true

	session := sessions.NewSession(&mockSession{
		save: func(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
			s.Values = make(map[interface{}]interface{})
			return nil
		},
	}, "session")

	authReq := &request.AuthReq{
		Request:  req,
		Locale:   h.Locale,
		Session:  session,
		Response: httptest.NewRecorder(),
		Data:     make(map[string]interface{}),
	}

	if user != nil {
		authReq.AuthUser = request.NewAuthUser(s, user, true, time.Duration(time.Hour))
	}

	return authReq
}

func prepareClientAuthService() *mockAuthService {
	authService := makeMockAuthService()
	return authService
}

func prepareClientAuthHandlers(authService authService, s *settings.Settings) *AuthHandlers {
	return &AuthHandlers{
		Log:         zap.NewNop(),
		Locale:      locale.Static(),
		AuthService: authService,
		Settings:    s,
	}
}
