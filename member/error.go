package member

import "errors"

//ErrRegisteredDataNotMap errors rasied when registered user data struct is not a map struct.
var ErrRegisteredDataNotMap = errors.New("registered user data is not a map")

//ErrAccountRegisterExists errors rasied when registered account ecists.
var ErrAccountRegisterExists = errors.New("account registered exists")

//ErrUserBanned errors rasied when user status is banned.
var ErrUserBanned = errors.New("user banned")

//ErrUserNotFound errors rasied when user is not found.
var ErrUserNotFound = errors.New("user not found")

//ErrFeatureNotSupported errors rasied when feature not supported by provider.
var ErrFeatureNotSupported = errors.New("feature not supported")

//ErrAccountKeywordNotRegistered errors rasied when account keyword not regietered.
var ErrAccountKeywordNotRegistered = errors.New("account keyword not regietered")

// ErrStatusNotSupport errors rasied when user status is not support by provider.
var ErrStatusNotSupport = errors.New("user status not  support")

//ErrPasswordNotChangeable errors raised when password provider not support change password.
var ErrPasswordNotChangeable = errors.New("password not changeable")
