module my.com/secrets/internal/auth/domain/test/e2e/hydra-kratos-login-consent

go 1.16

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/ory/hydra-client-go v1.7.4
	my.com/secrets/internal/auth/domain/-client-go v0.10.1
	github.com/ory/x v0.0.577
)