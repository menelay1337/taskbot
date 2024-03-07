package errors

func WrapIfErr(message string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(message, err)
}
