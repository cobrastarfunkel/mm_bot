package main

const (
	BOTNAME = "bobbot"
	TOKEN   = "j4m37kzow78zprw5bwgtysz89a"
	MMURI   = "http://localhost:8065"
	WSURI   = "ws://localhost:8065"
)

var logger = Logger{}

func main() {
	logger.init(DEBUG)
	logger.Info("Starting Bobbot")

	c := Conn{}
	c.init()

	b := MMBot{botName: BOTNAME, teamName: "botsample", conn: &c}
	b.init()
}
