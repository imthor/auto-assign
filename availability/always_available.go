package availability

type AlwaysAvailable struct{}

func (c *AlwaysAvailable) IsAvailable(username string) (bool, error) {
	return true, nil
}
