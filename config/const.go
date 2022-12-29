package config

/**
 * @author Rancho
 * @date 2022/12/29
 */

type Env string

const (
	EnvUnknown = ""
	EnvLocal   = "local"
	EnvGithub  = "github"
	EnvDev     = "dev"
	EnvProd    = "prod"
)

func (e Env) IsProd() bool {
	return string(e) == "prod"
}
