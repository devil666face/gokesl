package installer

type Redhat struct {
}

func NewRedhat() *Redhat {
	return &Redhat{}
}

func (r *Redhat) Install() error {
	return nil
}
