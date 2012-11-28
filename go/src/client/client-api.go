package client

type Client interface {
	Map(chunk string) string
	Reduce(keyValues map[string][]string) string
}

func Pack(key string, val string) string {
	return key + ", " + val + "\n"
}

func Unpack(keyVal string) string {
	return ""
}