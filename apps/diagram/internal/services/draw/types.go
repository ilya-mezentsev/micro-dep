package draw

type D2Client interface {
	Draw(script string) (string, error)
}
