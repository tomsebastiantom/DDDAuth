package service

import "my.com/secrets/internal/others/domain/translation/entity"

type Translator interface {
	Translate(translation entity.Translation) (entity.Translation, error)
}
