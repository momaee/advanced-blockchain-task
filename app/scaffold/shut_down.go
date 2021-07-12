package scaffold

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

func (s *skeleton) Close() {
	fmt.Println(aurora.Yellow("\n\n###############  Shutting Down ############### ").BgBlue())
	defer s.cancel()
	os.Exit(1)
}
