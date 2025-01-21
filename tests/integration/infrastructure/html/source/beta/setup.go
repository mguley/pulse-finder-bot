package beta

// SetupTestContainer initializes the TestContainer and handles cleanup.
func SetupTestContainer() *TestContainer {
	return NewTestContainer()
}
