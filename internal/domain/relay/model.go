package relay

import (
    validation "github.com/go-ozzo/ozzo-validation"
)
type JSON struct {
}

type Service struct {

}

type MySQL struct {
}

func (j *JSON) ToService() *Service {
    return &Service{
    }
}

func (s *Service) ToJSON() *JSON {
    return &JSON{
    }
}

func (p *MySQL) ToService() *Service {
    return &Service{
    }
}

func (s *Service) ToMySQL() *MySQL {
    return &MySQL{
    }
}

func (s *Service) Validate() error {
    return validation.ValidateStruct(s)
}
