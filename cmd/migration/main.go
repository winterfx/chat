package migration

import (
	"fmt"
	"os/exec"
)

func Start() {
	//run bash `sqlc generate`
	cmd := exec.Command("sqlc", "generate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running sqlc generate: %v\n", err)
		fmt.Printf("Output: %s\n", output)
		return
	}
	fmt.Println("sqlc generate success!")
}
