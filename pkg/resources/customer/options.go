package customer

// LoginOption is a functional option for the Login operation.
type LoginOption func(*LoginOptions)

// VerifyCodeOption is a functional option for the VerifyCode operation.
type VerifyCodeOption func(*VerifyCodeOptions)

// OAuthOption is a functional option for OAuth login operations (Apple, Google).
type OAuthOption func(*OAuthOptions)
